import {useEffect, useState} from 'react'
import mqtt from 'precompiled-mqtt'
import './StateReport.css'


const client = mqtt.connect("ws://DEPTH:9001")
client.subscribe("asol/#", console.log)

export default function StateReport() {
	const [stateReport, setStateReport] = useState("")
	const [connected, setConnected] = useState(false)

	useEffect(() => {
		client.on("message", (topic, payload) => {
			if (topic === "asol/state-report-json") {
				setStateReport(payload.toString())
			}
			console.log(topic, String(payload))
			console.log(client.connected)
		})
		client.on("connect", () => {setConnected(true)})
		client.on("close", () => {setConnected(false)})
	}, [])

	return (
		<>
		<h2>StateReport</h2>
		Connection: {String(connected)}
		<br/>
		<textarea id="stateReport" readOnly value={stateReport}></textarea>
		<button onClick={()=>{client.publish("test", "ahh")}}>Pub</button>
		<br/>
		</>
	)
}
	