import dayjs from "dayjs";
import { useEffect } from "react";
import { useClaim, useUnclaim } from "../ebs/api";
import { useGlobalState } from "../helpers/State";

export default function ConnectionManager() {
	const gs = useGlobalState();

	let connectionValid = false;
	if (gs.ebsState) {
		if (dayjs(gs.ebsState.ConnectedUserExpiryTimestamp).isAfter(dayjs())) {
			connectionValid = true;
		}
	}
	const connectedUser: string | undefined = connectionValid
		? gs.ebsState?.ConnectedUser?.OUID
		: undefined;

	const { mutate: claim, isPending: isClaiming } = useClaim();
	const { mutate: unclaim, isPending: isUnclaiming } = useUnclaim();

	const connectHandler = () => {
		claim();
	};

	const disconnectHandler = () => {
		unclaim();
	};

	const isThisUserConnected =
		connectedUser !== undefined && connectedUser === gs.auth?.userId;

	console.log("connectionValid", connectionValid);
	console.log("connectedUser", connectedUser);
	console.log("isThisUserConnected", isThisUserConnected);

	// regularly claim if this user is connected
	// useEffect(() => {
	// 	if (isThisUserConnected) {
	// 		const interval = setInterval(() => {
	// 			claim();
	// 		}, 5000);
	// 		return () => clearInterval(interval);
	// 	}
	// }, [isThisUserConnected, claim]);

	return (
		<div id="connection-area">
			{connectedUser ? (
				isThisUserConnected ? (
					<div className="connection-status">You are connected</div>
				) : (
					<div className="connection-status">A user is connected</div>
				)
			) : (
				<>
					<div className="connection-status">No connected user</div>
					<button
						type="button"
						onClick={connectHandler}
						// onKeyDown={connectHandler(auth)}
					>
						Connect
					</button>
					<button
						type="button"
						onClick={disconnectHandler}
						// onKeyDown={connectHandler(auth)}
					>
						Disconnect
					</button>
				</>
			)}
		</div>
	);
}
