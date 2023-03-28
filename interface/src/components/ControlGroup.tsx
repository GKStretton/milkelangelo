import { useState, useContext } from 'react';
import { SLEEP_TOPIC, TOPIC_STATE_REPORT_JSON, WAKE_TOPIC, SET_VALVE_TOPIC, SHUTDOWN_TOPIC, COLLECT_TOPIC, DISPENSE_TOPIC, FLUID_REQ_TOPIC, VALVE_DRAIN, VALVE_WATER, VALVE_MILK, VALVE_AIR } from '../util/topics'
import MqttContext from '../util/mqttContext'
import { ButtonGroup, Button, Typography, Slider, Box } from '@mui/material';

export default function ControlGroup() {
    const { client: c, messages } = useContext(MqttContext);
    const stateReportRaw = messages[TOPIC_STATE_REPORT_JSON];
    const stateReport = stateReportRaw && JSON.parse(stateReportRaw);

    const [dispenseVolume, setDispenseVolume] = useState(10.0);
    const [collectionVolume, setCollectionVolume] = useState(30.0);
    const [bulkFluidRequestVolume, setBulkFluidRequestVolume] = useState(200.0);
    // Create an array of marks with a 5µl interval
    const marks = Array.from({ length: 21 }, (_, i) => {
        return { value: i * 5, label: `${i * 5}µl` };
    });
    const marks_ml = Array.from({ length: 21 }, (_, i) => {
        return { value: i * 50, label: `${i * 50}ml` };
    });


    const isAwake = stateReport && stateReport.status !== "SLEEPING" && stateReport.status !== "E_STOP_ACTIVE";

    const noVials = 7
    const vials = new Array(noVials).fill(0).map((_, i) => noVials - i);

    const collecting = stateReport && stateReport.collection_request.completed === false;
    const collectingVial = collecting && stateReport.collection_request.vial_number;

    const bulkRequests = [
        {id: 1, name: "Milk", fluid_type: 3, open_drain: false},
        {id: 2, name: "Water", fluid_type: 2, open_drain: false},
        {id: 3, name: "Rinse", fluid_type: 2, open_drain: true},
        {id: 4, name: "Drain", fluid_type: 1, open_drain: false},
    ];

    const valves = [
        { id: VALVE_DRAIN, name: "DRAIN" },
        { id: VALVE_WATER, name: "WATER" },
        { id: VALVE_MILK, name: "MILK" },
        { id: VALVE_AIR, name: "AIR" }
    ];

    return (
        <>
        <Box
            display="flex"
            flexDirection="column"
            alignItems="center"
            justifyContent="center"
            m={4}
        >
            <Typography variant="h6">Basic</Typography>
            <ButtonGroup variant="outlined" aria-label="outlined button group" sx={{margin: 2}}>
                <Button variant="contained" disabled={isAwake} onClick={() => c?.publish(WAKE_TOPIC, "")}>Wake</Button>
                <Button disabled={!isAwake} onClick={() => c?.publish(SHUTDOWN_TOPIC, "")}>Shutdown</Button>
                <Button disabled={!isAwake} variant="contained" color="error" onClick={() => c?.publish(SLEEP_TOPIC, "")}>Kill</Button>
            </ButtonGroup>
            <Typography variant="h6">Collection</Typography>
            <Slider
                value={collectionVolume}
                onChange={(e, value) => typeof value === "number" ? setCollectionVolume(value): null}
                min={1}
                max={50} // Adjust the max value according to your requirements
                step={1}
                marks={marks}
                valueLabelDisplay="on"
                valueLabelFormat={(value) => `${value}µl`}
                aria-label="Dispense volume"
                sx={{margin: 5, width: "50%"}}
            />
            <ButtonGroup variant="outlined" aria-label="outlined button group" sx={{margin: 2}}>
                {vials.map((vial) => 
                    <Button
                        key={vial}
                        disabled={!isAwake || collecting}
                        variant={collectingVial === vial ? "contained": "outlined"}
                        onClick={() => c?.publish(COLLECT_TOPIC, `${vial.toString()},${collectionVolume}`)}
                    >
                        {vial}
                    </Button>
                )}
            </ButtonGroup>
            <Typography variant="h6">Dispense</Typography>
            <Slider
                value={dispenseVolume}
                onChange={(e, value) => typeof value === "number" ? setDispenseVolume(value): null}
                min={1}
                max={collectionVolume} // Adjust the max value according to your requirements
                step={1}
                marks={marks}
                valueLabelDisplay="on"
                valueLabelFormat={(value) => `${value}µl`}
                aria-label="Dispense volume"
                sx={{margin: 4, width: "50%"}}
            />
            <Button disabled={!isAwake || collecting} onClick={() => c?.publish(DISPENSE_TOPIC, dispenseVolume.toString())} sx={{"margin": 2}}>Dispense</Button>
            <Typography variant="h6">Bulk Fluid</Typography>
            <Slider
                value={bulkFluidRequestVolume}
                onChange={(e, value) => typeof value === "number" ? setBulkFluidRequestVolume(value): null}
                min={25}
                max={300} // Adjust the max value according to your requirements
                step={25}
                marks={marks_ml}
                valueLabelDisplay="on"
                valueLabelFormat={(value) => `${value}ml`}
                aria-label="Bulk Fluid Request volume"
                sx={{margin: 5, width: "50%"}}
            />
            <ButtonGroup variant="outlined" aria-label="outlined button group" sx={{margin: 2}}>
                {bulkRequests.map((request) => 
                    <Button
                        key={request.id}
                        disabled={!isAwake}
                        onClick={() => c?.publish(FLUID_REQ_TOPIC, `${request.fluid_type},${bulkFluidRequestVolume},${request.open_drain}`)}
                    >
                        {request.name}
                    </Button>
                )}
            </ButtonGroup>
            <Typography variant="h6">Valve Overrides (Open/Close)</Typography>
            <Box display="flex" flexDirection="row" alignItems="center" justifyContent="center" sx={{margin: 2}}>
                {valves.map(valve => (
                    <ButtonGroup key={valve.id} sx={{margin:1}}>
                        <Button disabled={!isAwake} onClick={() => c?.publish(SET_VALVE_TOPIC, `${valve.id},true`)}>{valve.name} (O)</Button>
                        <Button disabled={!isAwake} onClick={() => c?.publish(SET_VALVE_TOPIC, `${valve.id},false`)}>{valve.name} (C)</Button>
                    </ButtonGroup>
                ))}
            </Box>
        </Box>
        </>
    )
}
