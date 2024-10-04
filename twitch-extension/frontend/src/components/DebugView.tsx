import React, { useMemo, useState } from "react";

export default function DebugView({ robotState }: { robotState: undefined }) {
	const [show, setShow] = useState(true);

	const stateText = useMemo(
		() => JSON.stringify(robotState, null, " "),
		[robotState],
	);

	return (
		<>
			{show ? (
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
					left: show ? "0vw" : "0vw",
				}}
				onClick={() => {
					setShow(!show);
				}}
			>
				Dbg
			</div>
		</>
	);
}
