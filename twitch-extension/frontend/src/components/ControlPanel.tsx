import _ from "lodash";
import React from "react";
import { collectRequest, dispenseRequest } from "../ebs/api";
import { createCollectionVote } from "../ebs/helpers";

export default function ControlPanel({
	auth,
}: { auth: Twitch.ext.Authorized }) {
	const collectionHandler =
		(auth: Twitch.ext.Authorized, vialPos: number) => () => {
			if (!auth) return;
			collectRequest(auth, vialPos);
		};
	const dispenseHandler = (auth: Twitch.ext.Authorized) => () => {
		if (!auth) return;
		dispenseRequest(auth);
	};
	return (
		<div id="color-vote-area">
			{_.times(5, (i) => {
				const vialPos = i + 2;
				return (
					<div
						className="color-option"
						onClick={collectionHandler(auth, vialPos)}
						onKeyDown={collectionHandler(auth, vialPos)}
					>
						{vialPos}
					</div>
				);
			})}
			<div
				className="dispense-button"
				onClick={dispenseHandler(auth)}
				onKeyDown={dispenseHandler(auth)}
			>
				Dispense
			</div>
		</div>
	);
}
