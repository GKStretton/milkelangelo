import {
	QueryClient,
	QueryClientProvider,
	useQuery,
} from "@tanstack/react-query";
import _ from "lodash";
import { useEffect } from "react";
import { Toaster, toast } from "sonner";
import "./App.css";
import ConnectionManager from "./components/ConnectionManager";
import ControlPanel from "./components/ControlPanel";
import ControlView from "./components/ControlView";
import DebugView from "./components/DebugView";
import { setupApiAuth, useDirectEbsState } from "./ebs/api";
import { useGlobalState } from "./helpers/State";

function App() {
	const ext = window?.Twitch?.ext;
	const gs = useGlobalState();

	// only when testing
	const { data: directEbsState, error } = useDirectEbsState(gs.isLocalMode);

	useEffect(() => {
		if (gs.isLocalMode) {
			// use direct ebs state getter
			if (directEbsState) {
				gs.setEbsState(directEbsState);
			}
		}
	}, [gs.isLocalMode, directEbsState, gs.setEbsState]);

	useEffect(() => {
		if (error) {
			toast.error(`Failed to directly get EBS state: ${error.message}`);
		}
	}, [error]);

	useEffect(() => {
		if (ext.viewer.id === null) {
			toast.message("using local mode");
			gs.setLocalMode(true);
			gs.setAuth({
				userId: "local",
				channelId: "local",
				clientId: "local",
				token: "local",
				helixToken: "local",
			});
			return;
		}

		ext.onAuthorized((auth) => {
			console.log("got auth: ", auth);
			gs.setAuth(auth);
		});
		// todo: migrate from pubsub to eventsub?
		ext.listen("broadcast", (target, contentType, message) => {
			console.log("got broadcast: ", target, contentType, message);
			gs.setEbsState(JSON.parse(message));
		});
	}, [ext, gs.setAuth, gs.setEbsState, gs.setLocalMode]);

	useEffect(() => {
		if (gs.auth) {
			setupApiAuth(gs.auth);
		}
	}, [gs.auth]);

	return (
		<div className={`App ${gs.isLocalMode ? "mock-stream" : ""}`}>
			<Toaster richColors position="bottom-center" />
			<DebugView />
			<ConnectionManager />
			{gs.ebsState?.ConnectedUser && (
				<>
					<ControlView />
					<ControlPanel />
				</>
			)}
		</div>
	);
}

export default App;
