import { CssBaseline } from "@mui/material";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import "./App.css";
import KeyPressHandler from "./KeyPressHandler";
import CleaningPage from "./components/CleaningPage";
import ContentPage from "./components/ContentPage";
import DevPage from "./components/DevPage";
import { ErrorManager } from "./components/ErrorManager";
import Header from "./components/Header";
import MqttProvider from "./util/MqttProvider";

function App() {
	const mqtt_url =
		process.env.REACT_APP_MQTT_URL ?? `ws://${window.location.hostname}:9001`;
	console.log(mqtt_url);

	return (
		<>
			<CssBaseline />
			<ErrorManager>
				<MqttProvider url={mqtt_url}>
					<Router>
						<KeyPressHandler />
						<Header />
						<Routes>
							<Route path="/" element={<DevPage />} />
							<Route path="/cleaning" element={<CleaningPage />} />
							<Route path="/content" element={<ContentPage />} />
						</Routes>
					</Router>
				</MqttProvider>
			</ErrorManager>
		</>
	);
}

export default App;
