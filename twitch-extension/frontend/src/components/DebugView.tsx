import React, { useMemo, useState } from "react";
import { toast } from "sonner";
import { useGlobalState } from "../helpers/State";

export default function DebugView() {
	const gs = useGlobalState();

	const stateText = useMemo(
		() => JSON.stringify(gs.ebsState, null, " "),
		[gs.ebsState],
	);

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
				onClick={() => {
					gs.setDebugMode(!gs.isDebugMode);
					toast.info(`Debug mode ${gs.isDebugMode ? "off" : "on"}`);
				}}
			>
				Dbg
			</div>
		</>
	);
}
