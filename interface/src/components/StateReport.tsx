import {useEffect, useContext } from 'react'
import './StateReport.css'
import { TOPIC_SESSION_STATUS_GET, TOPIC_SESSION_STATUS_RESP_JSON, TOPIC_STATE_REPORT_JSON, TOPIC_STREAM_STATUS_GET, TOPIC_STREAM_STATUS_RESP_JSON } from '../topics_backend/topics_backend'
import { TOPIC_STATE_REPORT_REQUEST } from '../topics_firmware/topics_firmware'
import MqttContext from '../util/mqttContext'
import { Grid, Button } from '@mui/material'

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
			console.log("subsribed to state report", m);
		});
		c.subscribe(TOPIC_SESSION_STATUS_RESP_JSON, (m) => {
			console.log("subsribed to session status", m);
		});
		c.subscribe(TOPIC_STREAM_STATUS_RESP_JSON, (m) => {
			console.log("subsribed to stream status", m);
		});
		c.publish(TOPIC_STATE_REPORT_REQUEST, "");
		c.publish(TOPIC_SESSION_STATUS_GET, "");
		c.publish(TOPIC_STREAM_STATUS_GET, "");
	}, [c?.connected])

	return (
		<>
		<div style={{ display: 'flex', flexDirection: 'column', padding: "5px"}}>
			Connection: {String(connected)}
			<br/>
			<Grid container flexDirection="row">
				<Grid item>
					<h2>State Report</h2>
					<textarea id="stateReport" readOnly value={stateReport}></textarea>
					<Button variant="contained" sx={{width: "100px", margin: "5px"}} onClick={()=>{c?.publish(TOPIC_STATE_REPORT_REQUEST, "")}}>Request</Button>
				</Grid>
				<Grid item>
					<h2>Session Status</h2>
					<textarea id="sessionStatus" readOnly value={sessionStatus}></textarea>
					<Button variant="contained" sx={{width: "100px", margin: "5px"}} onClick={()=>{c?.publish(TOPIC_SESSION_STATUS_GET, "")}}>Request</Button>
				</Grid>
				<Grid item>
					<h2>Stream status</h2>
					<textarea id="streamStatus" readOnly value={streamStatus}></textarea>
					<Button variant="contained" sx={{width: "100px", margin: "5px"}} onClick={()=>{c?.publish(TOPIC_STREAM_STATUS_GET, "")}}>Request</Button>
				</Grid>
			</Grid>
			<br/>
		</div>
		</>
	)
}
	