import { Grid } from "@mui/material";
import ControlGroup from "./ControlGroup";
import StateReport from "./StateReport";
import TopCamPlayer from "./TopCamPlayer";
import VideoPlayer from "./VideoPlayer";

export default function DevPage() {
  const webrtc_url = process.env.REACT_APP_WEBRTC_URL ?? "ws://DEPTH:8889";
  return (
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
  );
}