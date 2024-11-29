import React, { useMemo, useState } from "react";
import { useGlobalState } from "../helpers/State";

export default function DebugView({ robotState }: { robotState: undefined }) {
	const stateText = useMemo(
		() => JSON.stringify(robotState, null, " "),
		[robotState],
	);

	const gs = useGlobalState();

	return (
		<>
			{gs.isDebugMode ? (
				<>
					<textarea
						readOnly={true}
						className="debug-text"
						value={`State text: ${stateText}`}
					/>
					<div className="border-area" />
					<div className="safe-area" />
				</>
			) : null}
			{/* biome-ignore lint/a11y/useKeyWithClickEvents: debug */}
			<div
				className="debug-toggle"
				style={{
					left: gs.isDebugMode ? "0vw" : "0vw",
				}}
				onClick={() => {
					gs.setDebugMode(!gs.isDebugMode);
				}}
			>
				Dbg
			</div>
		</>
	);
}
