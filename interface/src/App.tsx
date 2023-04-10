import StateReport from './components/StateReport'
import Header from './components/Header';
import './App.css';
import { CssBaseline, Grid, Toolbar, Container } from '@mui/material';
import { PrecisionManufacturing } from '@mui/icons-material';
import MqttProvider from './util/MqttProvider';
import ControlGroup from './components/ControlGroup';
import VideoPlayer from './components/VideoPlayer';
import TopCamPlayer from './components/TopCamPlayer';
import KeyPressHandler from './KeyPressHandler';

function App() {

  return (
    <>
      <CssBaseline />
      <MqttProvider url='ws://DEPTH:9001'>
        <KeyPressHandler/>
        <Header/>
        <Grid container
          direction="row"
          spacing={2}
          justifyContent="flex-start"
          alignItems="flex-start"
          padding={2}
        >
          <Grid item
            container
            direction="column"
            xs={4}>
            <StateReport/>
          </Grid>
          <Grid item container direction="column" xs={4} justifyContent="center" alignItems="center">
            <ControlGroup/>
          </Grid>
          <Grid item container direction="column" xs={3} justifyContent="flex-start">
            <Grid item xs={12} sm={6} justifyContent="center" alignItems="center">
              <TopCamPlayer/>
            </Grid>
            <Grid item xs={12} sm={6} justifyContent="center" alignItems="center">
              <VideoPlayer url="DEPTH:8889/front-cam-crop/" name="front" show={false}/>
            </Grid>
          </Grid>
        </Grid>
      </MqttProvider>
    </>
  );
}

export default App;
