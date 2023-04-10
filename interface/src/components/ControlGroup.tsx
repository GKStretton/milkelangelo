import { useState, useContext, useEffect } from 'react';
import { TOPIC_SHUTDOWN, TOPIC_SLEEP, TOPIC_WAKE, TOPIC_SET_VALVE, TOPIC_DISPENSE, TOPIC_FLUID, TOPIC_COLLECT, TOPIC_SET_IK_Z } from '../topics_firmware/topics_firmware';
import { TOPIC_SESSION_BEGIN, TOPIC_SESSION_END, TOPIC_SESSION_PAUSE, TOPIC_SESSION_RESUME, TOPIC_STATE_REPORT_JSON, TOPIC_STREAM_END, TOPIC_STREAM_START } from '../topics_backend/topics_backend';
import { SessionStatus, SolenoidValve, StateReport, Status, StreamStatus } from '../machinepb/machine_pb';
import MqttContext from '../util/mqttContext'
import { ButtonGroup, Button, Typography, Slider, Box, Tabs, Tab } from '@mui/material';
import { useSessionStatus, useStateReport, useStreamStatus } from '../util/hooks';

export default function ControlGroup() {
	const [tabValue, setTabValue] = useState(0);
	const handleChange = (event: React.SyntheticEvent, newValue: number) => {
		setTabValue(newValue);
	}

    const { client: c, messages } = useContext(MqttContext);
    const stateReport: StateReport | null = useStateReport();
    const sessionStatus: SessionStatus | null = useSessionStatus();
    const streamStatus: StreamStatus | null = useStreamStatus();

    const [dispenseVolume, setDispenseVolume] = useState(10.0);
    const [collectionVolume, setCollectionVolume] = useState(30.0);
    const [bulkFluidRequestVolume, setBulkFluidRequestVolume] = useState(200.0);
    const [zLevel, setZLevel] = useState(42);

    useEffect(() => {
        c?.publish(TOPIC_SET_IK_Z, zLevel.toString());
    });

    // Create an array of marks with a 5µl interval
    const marks = Array.from({ length: 21 }, (_, i) => {
        return { value: i * 10, label: `${i * 10}µl` };
    });
    const marks_ml = Array.from({ length: 21 }, (_, i) => {
        return { value: i * 50, label: `${i * 50}ml` };
    });
    const zMarks = Array.from({ length: 21 }, (_, i) => {
        return { value: i * 5, label: `${i * 5}mm` };
    });


    const isAwake: boolean = !!stateReport &&
        stateReport?.getStatus() !== Status.SLEEPING &&
        stateReport?.getStatus() !== Status.E_STOP_ACTIVE;

    const noVials = 7
    const vials = new Array(noVials).fill(0).map((_, i) => noVials - i);

    const collecting: boolean = !!stateReport && stateReport?.getCollectionRequest()?.getCompleted() === false;
    const collectingVial = collecting && stateReport?.getCollectionRequest()?.getVialNumber();

    const bulkRequests = [
        {id: 1, name: "Milk", fluid_type: 3, open_drain: false},
        {id: 2, name: "Water", fluid_type: 2, open_drain: false},
        {id: 3, name: "Rinse", fluid_type: 2, open_drain: true},
        {id: 4, name: "Drain", fluid_type: 1, open_drain: false},
    ];

    const valves = [
        { id: SolenoidValve.VALVE_DRAIN, name: "DRAIN" },
        { id: SolenoidValve.VALVE_WATER, name: "WATER" },
        { id: SolenoidValve.VALVE_MILK, name: "MILK" },
        { id: SolenoidValve.VALVE_AIR, name: "AIR" }
    ];

    const keyDownHandler = (event: KeyboardEvent) => {
        const key = event.key;

        const num = parseInt(key, 10);
        if (num >= 1 && num <= 7) {
            c?.publish(TOPIC_COLLECT, `${key},${collectionVolume}`);
            return;
        }

        switch (key) {
            case ' ':
                c?.publish(TOPIC_DISPENSE, dispenseVolume.toString());
                break;
            case 'w':
                c?.publish(TOPIC_WAKE, "");
                break;
            case 's':
                c?.publish(TOPIC_SHUTDOWN, "");
                break;
            case 'k':
                c?.publish(TOPIC_SLEEP, "");
                break;
        }
    };

    useEffect(() => {
        window.addEventListener('keydown', keyDownHandler);
        return () => {
            window.removeEventListener('keydown', keyDownHandler);
        };
    }, [c, dispenseVolume, collectionVolume]);

    return (
        <>
        <Tabs value={tabValue}
            onChange={handleChange}
        >
            <Tab label="Core" />
            <Tab label="Overrides" />
            <Tab label="Sessions" />
            <Tab label="Socials" />
        </Tabs>

        {tabValue === 0 &&
        <>
            <Typography variant="h6">Basic</Typography>
            <ButtonGroup size="small" variant="outlined" aria-label="outlined button group" sx={{margin: 1}}>
                <Button variant="contained" disabled={isAwake} onClick={() => c?.publish(TOPIC_WAKE, "")}>Wake</Button>
                <Button disabled={!isAwake} onClick={() => c?.publish(TOPIC_SHUTDOWN, "")}>Shutdown</Button>
                <Button disabled={!isAwake} variant="contained" color="error" onClick={() => c?.publish(TOPIC_SLEEP, "")}>Kill</Button>
            </ButtonGroup>
            <ButtonGroup size="small" variant="outlined" aria-label="outlined button group" sx={{margin: 1}}>
                <Button disabled={!streamStatus || streamStatus.getLive()}
                    onClick={() => c?.publish(TOPIC_STREAM_START, "")}
                >Start Stream</Button>
                <Button disabled={!streamStatus || !streamStatus.getLive()}
                    onClick={() => c?.publish(TOPIC_STREAM_END, "")}
                >End Stream</Button>
            </ButtonGroup>
            <ButtonGroup size="small" variant="outlined" aria-label="outlined button group" sx={{margin: 1}}>
                <Button disabled={!sessionStatus || !sessionStatus.getComplete() || sessionStatus.getPaused()}
                    onClick={() => c?.publish(TOPIC_SESSION_BEGIN, "")}
                > Begin Session</Button>
                <Button disabled={!sessionStatus || !sessionStatus.getComplete() || sessionStatus.getPaused()}
                    variant="contained" onClick={() => c?.publish(TOPIC_SESSION_BEGIN, "PRODUCTION")}>Begin Prod. Session</Button>
                <Button disabled={!sessionStatus || !sessionStatus.getPaused()}
                    onClick={() => c?.publish(TOPIC_SESSION_RESUME, "")}>Resume Session</Button>
            </ButtonGroup>
            <ButtonGroup size="small" variant="outlined" aria-label="outlined button group" sx={{margin: 1}}>
                <Button disabled={!sessionStatus || sessionStatus.getComplete() || sessionStatus.getPaused()}
                    onClick={() => c?.publish(TOPIC_SESSION_PAUSE, "")}>Pause Session</Button>
                <Button disabled={!sessionStatus || sessionStatus.getComplete()}
                    onClick={() => c?.publish(TOPIC_SESSION_END, "")}>End Session</Button>
            </ButtonGroup>

            <Typography variant="h6">Bulk Fluid</Typography>
            <Slider
                value={bulkFluidRequestVolume}
                onChange={(e, value) => typeof value === "number" ? setBulkFluidRequestVolume(value): null}
                min={25}
                max={300} // Adjust the max value according to your requirements
                step={25}
                marks={marks_ml}
                valueLabelDisplay="auto"
                valueLabelFormat={(value) => `${value}ml`}
                aria-label="Bulk Fluid Request volume"
                sx={{margin: 2, width: "50%"}}
            />
            <ButtonGroup variant="outlined" aria-label="outlined button group" sx={{margin: 2}}>
                {bulkRequests.map((request) => 
                    <Button
                        key={request.id}
                        disabled={!isAwake}
                        onClick={() => c?.publish(TOPIC_FLUID, `${request.fluid_type},${bulkFluidRequestVolume},${request.open_drain}`)}
                    >
                        {request.name}
                    </Button>
                )}
            </ButtonGroup>

            <Typography variant="h6">Collection</Typography>
            <Slider
                value={collectionVolume}
                onChange={(e, value) => typeof value === "number" ? setCollectionVolume(value): null}
                min={1}
                max={100} // Adjust the max value according to your requirements
                step={1}
                marks={marks}
                valueLabelDisplay="auto"
                valueLabelFormat={(value) => `${value}µl`}
                aria-label="Dispense volume"
                sx={{margin: 2, width: "50%"}}
            />
            <ButtonGroup variant="outlined" aria-label="outlined button group" sx={{margin: 2}}>
                {vials.map((vial) => 
                    <Button
                        key={vial}
                        disabled={!isAwake || collecting}
                        variant={collectingVial === vial ? "contained": "outlined"}
                        onClick={() => c?.publish(TOPIC_COLLECT, `${vial.toString()},${collectionVolume}`)}
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
                valueLabelDisplay="auto"
                valueLabelFormat={(value) => `${value}µl`}
                aria-label="Dispense volume"
                sx={{margin: 2, width: "50%"}}
            />
            <Button disabled={!isAwake || collecting} onClick={() => c?.publish(TOPIC_DISPENSE, dispenseVolume.toString())} sx={{"margin": 2}}>Dispense</Button>

            <Typography variant="h6">Z-Level</Typography>
            <Slider
                value={zLevel}
                onChange={(e, value) => typeof value === "number" ? setZLevel(value): null}
                min={35}
                max={55} // Adjust the max value according to your requirements
                step={1}
                marks={zMarks}
                valueLabelDisplay="auto"
                valueLabelFormat={(value) => `${value}mm`}
                aria-label="Z-Level"
                sx={{margin: 2, width: "50%"}}
            />
        </>
        }

        {tabValue === 1 &&
        <>
            <Typography variant="h6">Open Valves</Typography>
            <Box display="flex" flexDirection="row" alignItems="center" justifyContent="center" sx={{margin: 0.5}}>
                <ButtonGroup size="small" sx={{margin:0.5}}>
                    {valves.map(valve => (
                        <Button key={valve.id} disabled={!isAwake} onClick={() => c?.publish(TOPIC_SET_VALVE, `${valve.id},true`)}>{valve.name} (O)</Button>
                    ))}
                </ButtonGroup>
            </Box>
            <Typography variant="h6">Close Valves</Typography>
            <Box display="flex" flexDirection="row" alignItems="center" justifyContent="center" sx={{margin: 0.5}}>
                <ButtonGroup size="small" sx={{margin:0.5}}>
                    {valves.map(valve => (
                        <Button key={valve.id} disabled={!isAwake} onClick={() => c?.publish(TOPIC_SET_VALVE, `${valve.id},false`)}>{valve.name} (C)</Button>
                    ))}
                </ButtonGroup>
            </Box>
        </>}

        {tabValue === 2 &&
        <>
            <Typography>
                Not implemented: view all sessions, delete sessions, prune, generate post., etc.
            </Typography>
        </>}

        {tabValue === 3 &&
        <>
            <Typography>
                Not implemented: post to social media
            </Typography>
        </>}
        </>
    )
}
