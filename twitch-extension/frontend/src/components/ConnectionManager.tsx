import { useEffect } from "react";
import { useClaim } from "../ebs/api";
import { useGlobalState } from "../helpers/State";

export default function ConnectionManager() {
	const gs = useGlobalState();

	// todo: get from ebsState
	const connectedUser: string | undefined = undefined;

	const { mutate: claim, isPending: isClaiming } = useClaim();

	const connectHandler = () => () => {
		claim();
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
			{connectedUser ? (
				connectedUser !== undefined && connectedUser === gs.auth?.userId ? (
					<div className="connection-status">You are connected</div>
				) : (
					<div className="connection-status">A user is connected</div>
				)
			) : (
				<>
					<div className="connection-status">No connected user</div>
					<button
						type="button"
						onClick={connectHandler()}
						// onKeyDown={connectHandler(auth)}
					>
						Connect
					</button>
				</>
			)}
		</div>
	);
}
