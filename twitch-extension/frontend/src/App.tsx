import _ from "lodash";
import { useEffect, useMemo, useState } from "react";
import "./App.css";
import ControlPanel from "./components/ControlPanel";
import ControlView from "./components/ControlView";
import DebugView from "./components/DebugView";
import { createCollectionVote, createLocationVote } from "./ebs/helpers";
import useWindowDimensions from "./hooks/WindowSize";

function App() {
	const ext = window?.Twitch?.ext;
	const [auth, setAuth] = useState<Twitch.ext.Authorized>();
	const [robotState, setRobotState] = useState();

	useEffect(() => {
		if (!ext) {
			console.error("ext not defined, not running on twitch?");
			return;
		}
		ext.onAuthorized((auth) => {
			console.log("got auth: ", auth);
			setAuth(auth);
		});
		ext.listen("broadcast", (target, contentType, message) => {
			console.log("got broadcast: ", target, contentType, message);
			setRobotState(JSON.parse(message));
		});
	}, [ext]);

	return (
		<div className="App">
			<header className="App-header">
				{ext && auth ? (
					<>
						<DebugView robotState={robotState} />
						<ControlPanel auth={auth} />
						<ControlView auth={auth} />
					</>
				) : (
					<p style={{ color: "#ff00ff" }}>
						Error: could not get auth from twitch!
					</p>
				)}
			</header>
		</div>
	);
}

export default App;
