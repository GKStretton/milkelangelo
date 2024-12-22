import { useEffect } from "react";
import { useClaim, useUnclaim } from "../ebs/api";
import { useGlobalState } from "../helpers/State";

export default function ConnectionManager() {
	const gs = useGlobalState();

	const haveGooState: boolean = !!gs.ebsState?.GooState;

	const connectedUser: string | undefined = gs.ebsState?.ConnectedUser?.OUID;

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

	// regularly claim if this user is connected
	useEffect(() => {
		if (isThisUserConnected) {
			const interval = setInterval(() => {
				claim();
			}, 5000);
			return () => clearInterval(interval);
		}
	}, [isThisUserConnected, claim]);

	return (
		<div id="connection-area">
			{haveGooState ? (
				<div className="connection-status green">Goo State available</div>
			) : (
				<div className="connection-status red">Goo State unavailable</div>
			)}
			{isThisUserConnected ? (
				<>
					<div className="connection-status green">You are connected</div>
					<button
						type="button"
						onClick={disconnectHandler}
						// onKeyDown={connectHandler(auth)}
					>
						Disconnect
					</button>
				</>
			) : null}
			{!isThisUserConnected && connectedUser !== undefined ? (
				<div className="connection-status amber">A user is connected</div>
			) : null}

			{connectedUser === undefined ? (
				<>
					<div className="connection-status red">No connected user</div>
					<button
						type="button"
						onClick={connectHandler}
						// onKeyDown={connectHandler(auth)}
					>
						Connect
					</button>
				</>
			) : null}
		</div>
	);
}
