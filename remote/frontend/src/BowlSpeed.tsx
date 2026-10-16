import { useEffect, useState } from "react";

// matches control.MaxBowlSpeed in goo
const MAX_SPEED = 1000;
// wait for the slider to settle before sending
const SEND_DELAY_MS = 300;
// show the sent speed until the robot reports it, or give up after this
const CONFIRM_TIMEOUT_MS = 2000;

interface BowlSpeedProps {
	// last requested speed, steps/sec, negative for reverse
	speed: number;
	disabled: boolean;
	onSet: (speed: number) => void;
}

export default function BowlSpeed({ speed, disabled, onSet }: BowlSpeedProps) {
	const [draft, setDraft] = useState<number | null>(null);
	const [sent, setSent] = useState(false);

	useEffect(() => {
		if (draft === null) return;
		if (sent) {
			const t = window.setTimeout(() => setDraft(null), CONFIRM_TIMEOUT_MS);
			return () => window.clearTimeout(t);
		}
		const t = window.setTimeout(() => {
			onSet(draft);
			setSent(true);
		}, SEND_DELAY_MS);
		return () => window.clearTimeout(t);
	}, [draft, sent]);

	// the robot reported a new speed
	useEffect(() => setDraft(null), [speed]);

	const change = (v: number) => {
		setSent(false);
		setDraft(v);
	};

	const value = draft ?? speed;

	return (
		<div className="group">
			<h2>Bowl</h2>
			<p className="readout">
				{value === 0 ? (
					<span className="muted">Stopped</span>
				) : (
					<>
						{value > 0 ? "Forward" : "Reverse"}
						<span className="muted"> · {Math.abs(value)} steps/s</span>
					</>
				)}
			</p>
			<div className="bowl-speed">
				<input
					type="range"
					min={-MAX_SPEED}
					max={MAX_SPEED}
					step={50}
					value={value}
					disabled={disabled}
					onChange={(e) => change(Number(e.target.value))}
					aria-label="Bowl speed, steps per second"
				/>
				<button
					type="button"
					className="quiet"
					disabled={disabled || value === 0}
					onClick={() => change(0)}
				>
					Stop
				</button>
			</div>
		</div>
	);
}
