import { useEffect, useContext, useState } from "react";
import "./StateReport.css";
import {
  TOPIC_SESSION_STATUS_GET,
  TOPIC_SESSION_STATUS_RESP_JSON,
  TOPIC_STATE_REPORT_JSON,
  TOPIC_STREAM_STATUS_GET,
  TOPIC_STREAM_STATUS_RESP_JSON,
} from "../topics_backend/topics_backend";
import { TOPIC_STATE_REPORT_REQUEST } from "../topics_firmware/topics_firmware";
import MqttContext from "../util/mqttContext";
import { Button, Tabs, Tab } from "@mui/material";
import { useProtoTopic } from "../util/hooks";

export default function StateReport() {
  const { client: c } = useContext(MqttContext);
  const connected = c?.connected;

  const stateReport = useProtoTopic(TOPIC_STATE_REPORT_JSON);
  const sessionStatus = useProtoTopic(TOPIC_SESSION_STATUS_RESP_JSON);
  const streamStatus = useProtoTopic(TOPIC_STREAM_STATUS_RESP_JSON);

  useEffect(() => {
    if (!c || !c.connected) {
      return;
    }
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
