import StateReport from './components/StateReport'
import './App.css';
import VerticalSlider from './components/VerticalSlider';
import { Typography, AppBar, Button, Card, CardActions, CardContent, CardMedia, CssBaseline, Grid, Toolbar, Container } from '@mui/material';
import { Castle, PrecisionManufacturing } from '@mui/icons-material';
import MqttProvider from './util/MqttProvider';
import ControlGroup from './components/ControlGroup';

function App() {
  return (
    <>
      <CssBaseline />
      <AppBar position="relative">
        <Toolbar>
          <PrecisionManufacturing/>
          <Typography variant="h6">
            A Study of Light
          </Typography>
        </Toolbar>
      </AppBar>
      <MqttProvider url='ws://DEPTH:9001'>
        <Container>
          <ControlGroup/>
          <StateReport/>
          <VerticalSlider/>
        </Container>
      </MqttProvider>
    </>
  );
}

export default App;
