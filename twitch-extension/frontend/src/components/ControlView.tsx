import { setUncaughtExceptionCaptureCallback } from "process";
import React, { useState } from "react";
import { useGoTo } from "../ebs/api";
import { useGlobalState } from "../helpers/State";
import { Status } from "../types";
import "./ControlView.css";

export default function ControlView() {
	const gs = useGlobalState();

	const isGoToEnabled =
		gs.ebsState?.GooState?.Status === Status.StatusDecidingDispense;
	// const isGoToEnabled =
	// 	gs.ebsState?.GooState?.DispenseState?.Completed === false || // placement in progress
	// 	gs.ebsState?.GooState?.CollectionState?.Completed === false; // collection in progress

	// get x and y from state
	const x = gs.ebsState?.GooState?.X ?? 0;
	const y = gs.ebsState?.GooState?.Y ?? 0;

	const perc_x = (x / 2 + 0.5) * 100;
	const perc_y = (0.5 - y / 2) * 100;

	const { mutate: goTo, isPending: isGoingTo } = useGoTo();

	const locationVoteHandler = (e: React.MouseEvent<HTMLElement>) => {
		if (!isGoToEnabled) {
			return;
		}

		const target = e.target as HTMLElement;
		const bounds = target.getBoundingClientRect();
		console.log("bounds", bounds);

		const raw_x = e.clientX - bounds.left;
		const raw_y = e.clientY - bounds.top;

		const perc_x = (raw_x / bounds.width) * 100;
		const perc_y = (raw_y / bounds.height) * 100;

		const x = (perc_x / 100 - 0.5) * 2;
		const y = -(perc_y / 100 - 0.5) * 2;

		goTo({ x, y });
	};

	return (
		<>
			<div
				className={`canvas ${gs.isDebugMode ? "debug" : ""} ${
					isGoToEnabled ? "" : "disable"
				}`}
				onClick={locationVoteHandler}
				onKeyDown={() => {
					console.log("key down not supported on canvas yet");
				}}
			>
				<div
					className="cursor"
					style={{ left: `${perc_x}%`, top: `${perc_y}%` }}
				/>
			</div>
			{gs.isDebugMode && (
				<div className="control-view-debug">
					<span>{`coords: ${x.toFixed(2)}, ${y.toFixed(2)}`}</span>
					<span>{`%: ${perc_x.toFixed(1)}%, ${perc_y.toFixed(1)}%`}</span>
				</div>
			)}
		</>
	);
}
