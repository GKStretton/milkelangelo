import React, { useState } from "react";
import { useGoTo } from "../ebs/api";
import { useGlobalState } from "../helpers/State";
import "./ControlView.css";

export default function ControlView() {
	const gs = useGlobalState();

	const [{ x, y }, setRawCoords] = useState({ x: 0, y: 0 });

	const { mutate: goTo, isPending: isGoingTo } = useGoTo();

	const locationVoteHandler = (e: React.MouseEvent<HTMLElement>) => {
		const target = e.target as HTMLElement;
		const bounds = target.getBoundingClientRect();
		const x = e.clientX - bounds.left;
		const y = e.clientY - bounds.top;
		setRawCoords({ x: x, y: y });
		const xMod = x / 400.0;
		const yMod = y / 400.0;

		goTo({ x: xMod, y: yMod });
	};

	return (
		<div
			className={`canvas ${gs.isDebugMode ? "debug" : ""}`}
			onClick={locationVoteHandler}
			onKeyDown={() => {
				console.log("key down not supported on canvas yet");
			}}
		>
			<div className="cursor" style={{ left: x, top: y }} />
		</div>
	);
}
