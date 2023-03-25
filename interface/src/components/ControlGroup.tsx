import { useState, useContext, useEffect } from 'react';
import { SLEEP_TOPIC, TOPIC_STATE_REPORT_JSON, WAKE_TOPIC, SHUTDOWN_TOPIC, COLLECT_TOPIC } from '../util/topics'
import MqttContext from '../util/mqttContext'
import { Grid, ButtonGroup, Button } from '@mui/material';

export default function ControlGroup() {
    const { client: c, messages } = useContext(MqttContext);
    const stateReportRaw = messages[TOPIC_STATE_REPORT_JSON];
    const stateReport = stateReportRaw && JSON.parse(stateReportRaw);

    const isAwake = stateReport && stateReport.status !== "SLEEPING" && stateReport.status !== "E_STOP_ACTIVE";

    const noVials = 7
    const vials = new Array(noVials).fill(0).map((_, i) => noVials - i);

    const collecting = stateReport && stateReport.collection_request.completed === false;
    const collectingVial = collecting && stateReport.collection_request.vial_number;
    const collectionVolume = 30.0;

    return (
        <>
        <ButtonGroup variant="outlined" aria-label="outlined button group" sx={{margin: 1}}>
            <Button variant="contained" disabled={isAwake} onClick={() => c?.publish(WAKE_TOPIC, "")}>Wake</Button>
            <Button disabled={!isAwake} onClick={() => c?.publish(SHUTDOWN_TOPIC, "")}>Shutdown</Button>
            <Button disabled={!isAwake} variant="contained" color="error" onClick={() => c?.publish(SLEEP_TOPIC, "")}>Kill</Button>
        </ButtonGroup>
        <br/>
        Collection:
        <br/>
        <ButtonGroup variant="outlined" aria-label="outlined button group" sx={{margin: 1}}>
            {vials.map((vial) => 
                <Button
                    disabled={!isAwake || collecting}
                    variant={collectingVial == vial ? "contained": "outlined"}
                    onClick={() => c?.publish(COLLECT_TOPIC, `${vial.toString()},${collectionVolume}`)}
                >
                    {vial}
                </Button>
            )}
        </ButtonGroup>
        </>
    )
}
