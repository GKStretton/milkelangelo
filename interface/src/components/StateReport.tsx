import { useEffect, useContext, useState } from "react";
import "./StateReport.css";
import {
  TOPIC_SESSION_STATUS_GET,
  TOPIC_SESSION_STATUS_RESP_JSON,
  TOPIC_SESSION_STATUS_RESP_RAW,
  TOPIC_STATE_REPORT_JSON,
  TOPIC_STREAM_STATUS_GET,
  TOPIC_STREAM_STATUS_RESP_JSON,
  TOPIC_STREAM_STATUS_RESP_RAW,
} from "../topics_backend/topics_backend";
import { TOPIC_STATE_REPORT_RAW, TOPIC_STATE_REPORT_REQUEST } from "../topics_firmware/topics_firmware";
import MqttContext from "../util/mqttContext";
import { Grid, Button, Tabs, Tab } from "@mui/material";

export default function StateReport() {
  const { client: c, messages } = useContext(MqttContext);
  const stateReport = messages[TOPIC_STATE_REPORT_JSON];
  const sessionStatus = messages[TOPIC_SESSION_STATUS_RESP_JSON];
  const streamStatus = messages[TOPIC_STREAM_STATUS_RESP_JSON];
  const connected = c?.connected;

  useEffect(() => {
    if (!c || !c.connected) {
      return;
    }
    // todo: sub to the raw for when parsing individual fields
    c.subscribe(TOPIC_STATE_REPORT_JSON, (m) => {
      console.log("subscribed to state report", m);
    });
    c.subscribe(TOPIC_STATE_REPORT_RAW, (m) => {
      console.log("subscribed to state report", m);
    });
    c.subscribe(TOPIC_SESSION_STATUS_RESP_JSON, (m) => {
      console.log("subscribed to session status", m);
    });
    c.subscribe(TOPIC_SESSION_STATUS_RESP_RAW, (m) => {
      console.log("subscribed to session status", m);
    });
    c.subscribe(TOPIC_STREAM_STATUS_RESP_JSON, (m) => {
      console.log("subscribed to stream status", m);
    });
    c.subscribe(TOPIC_STREAM_STATUS_RESP_RAW, (m) => {
      console.log("subscribed to stream status", m);
    });
    c.publish(TOPIC_STATE_REPORT_REQUEST, "");
    c.publish(TOPIC_SESSION_STATUS_GET, "");
    c.publish(TOPIC_STREAM_STATUS_GET, "");
  }, [c?.connected]);

  const [tabValue, setTabValue] = useState(0);
  const handleChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  let selectedMessage: string = "";
  let selectedTopic: string = "";
  let id: string = "";
  switch (tabValue) {
    case 0:
      selectedMessage = stateReport?.toString() ?? "";
      selectedTopic = TOPIC_STATE_REPORT_REQUEST;
      id = "stateReport";
      break;
    case 1:
      selectedMessage = sessionStatus?.toString() ?? "";
      selectedTopic = TOPIC_SESSION_STATUS_GET;
      id = "sessionStatus";
      break;
    case 2:
      selectedMessage = streamStatus?.toString() ?? "";
      selectedTopic = TOPIC_STREAM_STATUS_GET;
      id = "streamStatus";
      break;
  }

  return (
    <>
      Connection: {String(connected)}
      <br />
      <Tabs value={tabValue} onChange={handleChange} variant="scrollable" scrollButtons="auto">
        <Tab label="S.R" />
        <Tab label="Session" />
        <Tab label="Stream" />
        <Tab label="F.W Logs" />
        <Tab label="B.E Logs" />
      </Tabs>
      {tabValue <= 2 && (
        <>
          <Button
            variant="contained"
            sx={{ width: "100px", margin: "5px" }}
            onClick={() => {
              c?.publish(selectedTopic, "");
            }}
          >
            Request
          </Button>
          <textarea id={id} readOnly value={selectedMessage}></textarea>
        </>
      )}
      {tabValue >= 3 && <>not implemented</>}
    </>
  );
}
