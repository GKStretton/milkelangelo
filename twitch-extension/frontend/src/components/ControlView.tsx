import React, { useState } from "react";
import { goToRequest } from "../ebs/api";
import { createLocationVote } from "../ebs/helpers";

export default function ControlView({ auth }: { auth: Twitch.ext.Authorized }) {
	const [{ x, y }, setCoords] = useState({ x: 0, y: 0 });

	const locationVoteHandler = (e: React.MouseEvent<HTMLElement>) => {
		const target = e.target as HTMLElement;
		const bounds = target.getBoundingClientRect();
		const x = e.clientX - bounds.left;
		const y = e.clientY - bounds.top;
		setCoords({ x: x, y: y });
		const xMod = x / 100.0;
		const yMod = y / 100.0;

		if (!auth) return;
		goToRequest(auth, xMod, yMod);
	};

	return (
		<div
			className="canvas"
			onClick={locationVoteHandler}
			onKeyDown={() => {
				console.log("key down not supported on canvas yet");
			}}
		>
			<div className="cursor" style={{ left: x, top: y }} />
		</div>
	);
}
