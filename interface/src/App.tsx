import StateReport from "./components/StateReport";
import Header from "./components/Header";
import "./App.css";
import { CssBaseline, Grid, Toolbar, Container } from "@mui/material";
import { PrecisionManufacturing } from "@mui/icons-material";
import MqttProvider from "./util/MqttProvider";
import ControlGroup from "./components/ControlGroup";
import VideoPlayer from "./components/VideoPlayer";
import TopCamPlayer from "./components/TopCamPlayer";
import KeyPressHandler from "./KeyPressHandler";
import { ErrorManager } from "./components/ErrorManager";

function App() {
  const mqtt_url = process.env.REACT_APP_MQTT_URL ?? "ws://DEPTH:9001";
  const webrtc_url = process.env.REACT_APP_WEBRTC_URL ?? "ws://DEPTH:8889";
  console.log(mqtt_url);

  return (
    <>
      <CssBaseline />
      <ErrorManager>
        <MqttProvider url={mqtt_url}>
          <KeyPressHandler />
          <Header />
          <Grid container direction="row" spacing={2} justifyContent="flex-start" alignItems="flex-start" padding={2}>
            <Grid item container direction="column" xs={4}>
              <StateReport />
            </Grid>
            <Grid item container direction="column" xs={4} justifyContent="center" alignItems="center">
              <ControlGroup />
            </Grid>
            <Grid item container direction="column" xs={3} justifyContent="flex-start">
              <Grid item xs={12} sm={6} justifyContent="center" alignItems="center">
                <TopCamPlayer url={webrtc_url} />
              </Grid>
              <Grid item xs={12} sm={6} justifyContent="center" alignItems="center">
                <VideoPlayer url={`${webrtc_url}/front-cam-crop/`} name="front" />
              </Grid>
            </Grid>
          </Grid>
        </MqttProvider>
      </ErrorManager>
    </>
  );
}

export default App;
