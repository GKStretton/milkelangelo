import _ from "lodash";
import { useEffect, useMemo, useState } from "react";
import "./App.css";
import ControlPanel from "./components/ControlPanel";
import ControlView from "./components/ControlView";
import DebugView from "./components/DebugView";
import { Coords } from "./types";

function App() {
	const ext = window?.Twitch?.ext;
	const [auth, setAuth] = useState<Twitch.ext.Authorized>();
	const [authDisabled, setAuthDisabled] = useState(false);
	const [robotState, setRobotState] = useState();
	const [coords, setCoords] = useState<Coords>({ x: 0, y: 0 });

	useEffect(() => {
		if (ext.viewer.id === null) {
			console.log("disabling auth");
			setAuthDisabled(true);
			return;
		}

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
			{ext && (auth || authDisabled) ? (
				<>
					<DebugView robotState={robotState} />
					<ControlView auth={auth} coords={coords} setCoords={setCoords} />
					<ControlPanel auth={auth} coords={coords} />
				</>
			) : (
				<p style={{ color: "#ff00ff" }}>
					Error: could not get auth from twitch!
				</p>
			)}
		</div>
	);
}

export default App;
