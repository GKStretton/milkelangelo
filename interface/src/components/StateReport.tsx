import {useEffect, useContext } from 'react'
import './StateReport.css'
import { TOPIC_STATE_REPORT_JSON } from '../topics_backend/topics_backend'
import { TOPIC_STATE_REPORT_REQUEST } from '../topics_firmware/topics_firmware'
import MqttContext from '../util/mqttContext'
import { Grid, Button } from '@mui/material'

export default function StateReport() {
	const { client: c, messages } = useContext(MqttContext);
	const stateReport = messages[TOPIC_STATE_REPORT_JSON];
	const connected = c?.connected;

	useEffect(() => {
		if (!c || !c.connected) {
			return;
		}
		c.subscribe(TOPIC_STATE_REPORT_JSON, (m) => {
			console.log("subsribed to state report", m);
		});
		c.publish(TOPIC_STATE_REPORT_REQUEST, "");
	}, [c?.connected])

	return (
		<>
		<div style={{ display: 'flex', flexDirection: 'column', padding: "5px"}}>
			<h2>State Report</h2>
			Connection: {String(connected)}
			<br/>
			<textarea id="stateReport" readOnly value={stateReport}></textarea>
			<Button variant="contained" sx={{width: "100px", margin: "5px"}} onClick={()=>{c?.publish("mega/req/state-report", "")}}>Pub</Button>
			<br/>
		</div>
		</>
	)
}
	