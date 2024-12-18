import { connected } from "process";
import { useEffect } from "react";
import { claimRequest } from "../ebs/api";

export default function ConnectionManager({
	auth,
	ebsState,
}: {
	auth: Twitch.ext.Authorized | undefined;
	ebsState: undefined;
}) {
	// todo: get from ebsState
	const connectedUser: string | undefined = undefined;

	const connectHandler = (auth: Twitch.ext.Authorized | undefined) => () => {
		console.log("connect");
		claimRequest(auth);
	};

	const isThisUserConnected =
		connectedUser !== undefined && connectedUser === auth?.userId;

	// regularly claim if this user is connected
	useEffect(() => {
		if (isThisUserConnected) {
			const interval = setInterval(() => {
				claimRequest(auth);
			}, 5000);
			return () => clearInterval(interval);
		}
	}, [isThisUserConnected, auth]);

	return (
		<div id="connection-area">
			{connectedUser ? (
				connectedUser !== undefined && connectedUser === auth?.userId ? (
					<div className="connection-status">You are connected</div>
				) : (
					<div className="connection-status">A user is connected</div>
				)
			) : (
				<>
					<div className="connection-status">No connected user</div>
					<button
						type="button"
						onClick={connectHandler(auth)}
						// onKeyDown={connectHandler(auth)}
					>
						Connect
					</button>
				</>
			)}
		</div>
	);
}
