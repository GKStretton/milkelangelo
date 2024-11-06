import React, { useState } from "react";
import { goToRequest } from "../ebs/api";
import { createLocationVote } from "../ebs/helpers";
import { Coords } from "../types";

export default function ControlView({
	auth,
	coords,
	setCoords,
}: {
	auth: Twitch.ext.Authorized;
	coords: Coords;
	setCoords: (coords: Coords) => void;
}) {
	const [{ x, y }, setRawCoords] = useState({ x: 0, y: 0 });

	const locationVoteHandler = (e: React.MouseEvent<HTMLElement>) => {
		const target = e.target as HTMLElement;
		const bounds = target.getBoundingClientRect();
		const x = e.clientX - bounds.left;
		const y = e.clientY - bounds.top;
		setRawCoords({ x: x, y: y });
		const xMod = x / 400.0;
		const yMod = y / 400.0;

		setCoords({ x: xMod, y: yMod });

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
