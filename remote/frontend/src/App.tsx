import { useEffect } from "react";
import Bowl from "./Bowl";
import FrontView from "./FrontView";
import type { GooState, VialProfile } from "./types";
import { useControl } from "./useControl";
import { useRemoteState } from "./useRemoteState";

interface Vial {
	position: number;
	profile: VialProfile;
}

function vialsOf(gs: GooState | undefined): Vial[] {
	return Object.entries(gs?.VialProfiles ?? {})
		.filter((e): e is [string, VialProfile] => e[1] !== null)
		.map(([pos, profile]) => ({ position: Number(pos), profile }))
		.sort((a, b) => a.position - b.position);
}

export default function App() {
	const { state, connected } = useRemoteState();
	const control = useControl();
	const gs = state?.GooState;

	const vials = vialsOf(gs);
	const holding =
		gs?.DispenseState && !gs.DispenseState.Completed ? gs.DispenseState : null;
	const heldVial = vials.find((v) => v.position === holding?.VialNumber);

	const live = connected && !!gs;
	const enabled = live && gs.ControlEnabled;
	const asleep = live && gs.Status === "sleeping";
	const mine = enabled && control.hasControl;
	const canCollect = mine && gs.WaitingForCollection && !control.pending;
	const canDispense = mine && gs.WaitingForDispense && !control.pending;

	// Space drops, number keys pick a vial
	useEffect(() => {
		const onKey = (e: KeyboardEvent) => {
			if (e.target instanceof HTMLButtonElement && e.key === " ") return;
			if (e.key === " " && canDispense && gs) {
				e.preventDefault();
				control.dispense(gs.X, gs.Y);
			}
			const vial = vials.find((v) => String(v.position) === e.key);
			if (vial && canCollect) control.collect(vial.position, vial.profile.Name);
		};
		window.addEventListener("keydown", onKey);
		return () => window.removeEventListener("keydown", onKey);
	});

	let status: { tone: "off" | "wait" | "go"; text: string };
	if (!connected) status = { tone: "off", text: "Can't reach the robot. Retrying…" };
	else if (!gs) status = { tone: "wait", text: "Connecting…" };
	else if (!gs.ControlEnabled) status = { tone: "off", text: "Remote control is switched off" };
	else if (asleep) status = { tone: "off", text: "The robot is asleep" };
	else if (!control.hasControl && state.Claimed)
		status = { tone: "wait", text: "Someone else is painting right now" };
	else if (!control.hasControl) status = { tone: "go", text: "The robot is free" };
	else if (gs.WaitingForCollection) status = { tone: "go", text: "Pick a colour" };
	else if (gs.WaitingForDispense)
		status = { tone: "go", text: "Click the bowl to aim, then drop" };
	else status = { tone: "wait", text: "Working…" };

	return (
		<div className="page">
			<header className="bar">
				<h1>Milkelangelo</h1>
				<span className={`pill ${live ? "live" : ""}`}>{live ? "Live" : "Offline"}</span>
			</header>

			<main className="layout">
				<section className="stage">
					<Bowl
						x={gs?.X ?? 0}
						y={gs?.Y ?? 0}
						canAim={canDispense}
						onAim={control.goTo}
					/>
					<FrontView />
				</section>

				<aside className="panel">
					<div className={`status ${status.tone}`} aria-live="polite">
						{status.text}
					</div>

					{enabled && !asleep && (
						<div className="claim">
							{control.hasControl ? (
								<button type="button" className="quiet" onClick={control.release}>
									Give up control
								</button>
							) : (
								<button
									type="button"
									className="primary"
									onClick={control.claim}
									disabled={state?.Claimed || control.pending}
								>
									Take control
								</button>
							)}
						</div>
					)}

					<div className="group">
						<h2>Colours</h2>
						{vials.length === 0 ? (
							<p className="muted">No colours loaded</p>
						) : (
							<div className="vials">
								{vials.map((v) => (
									<button
										type="button"
										key={v.position}
										className={`vial ${heldVial === v ? "held" : ""}`}
										disabled={!canCollect}
										onClick={() => control.collect(v.position, v.profile.Name)}
									>
										<span
											className="swatch"
											style={{ background: v.profile.Colour || "currentColor" }}
										/>
										<span className="vial-name">{v.profile.Name}</span>
										<kbd>{v.position}</kbd>
									</button>
								))}
							</div>
						)}
					</div>

					<div className="group">
						<h2>Pipette</h2>
						<p className="readout">
							{holding ? (
								<>
									{heldVial?.profile.Name ?? `Vial ${holding.VialNumber}`}
									<span className="muted">
										{" "}
										· {holding.VolumeRemainingUl.toFixed(0)} µl left
									</span>
								</>
							) : (
								<span className="muted">Empty</span>
							)}
						</p>
						<button
							type="button"
							className="primary drop"
							disabled={!canDispense}
							onClick={() => gs && control.dispense(gs.X, gs.Y)}
						>
							Drop <kbd>Space</kbd>
						</button>
					</div>

					{control.notice && (
						<div className={`notice ${control.notice.tone}`} role="status">
							{control.notice.text}
						</div>
					)}
				</aside>
			</main>
		</div>
	);
}
