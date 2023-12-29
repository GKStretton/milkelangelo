import Header from "./components/Header";
import "./App.css";
import { CssBaseline } from "@mui/material";
import MqttProvider from "./util/MqttProvider";
import KeyPressHandler from "./KeyPressHandler";
import { ErrorManager } from "./components/ErrorManager";
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import DevPage from "./components/DevPage";
import CleaningPage from "./components/CleaningPage";
import ContentPage from "./components/ContentPage";

function App() {
  const mqtt_url = process.env.REACT_APP_MQTT_URL ?? "ws://DEPTH:9001";
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
              <Route path='/' element={<DevPage/>} />
              <Route path='/cleaning' element={ <CleaningPage/> } />
              <Route path='/content' element={ <ContentPage/> } />
            </Routes>
          </Router>
        </MqttProvider>
      </ErrorManager>
    </>
  );
}

export default App;
