import _ from "lodash";
import React from "react";
import { createCollectionVote } from "../ebs/helpers";

export default function ControlPanel({
	auth,
}: { auth: Twitch.ext.Authorized }) {
	const collectionVoteHandler =
		(auth: Twitch.ext.Authorized, i: number) => () => {
			if (!auth) return;
			console.log("voting collection for ", 6 - i);
			fetch("http://localhost:8080/vote", {
				method: "POST",
				body: JSON.stringify(createCollectionVote(6 - i)),
				headers: {
					Authorization: `Bearer ${auth.token}`,
					"X-Twitch-Extension-Client-Id": auth.clientId,
				},
			}).catch((e) => console.error(e));
		};
	return (
		<div id="color-vote-area">
			{_.times(5, (i) => (
				<div
					className="color-option"
					onClick={collectionVoteHandler(auth, i)}
					onKeyDown={collectionVoteHandler(auth, i)}
				/>
			))}
		</div>
	);
}
