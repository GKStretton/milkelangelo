import { useState, useContext, useEffect } from 'react';
import { TOPIC_SHUTDOWN, TOPIC_SLEEP, TOPIC_WAKE, TOPIC_SET_VALVE, TOPIC_DISPENSE, TOPIC_FLUID, TOPIC_COLLECT, TOPIC_SET_IK_Z } from '../topics_firmware/topics_firmware';
import { SessionStatus, SolenoidValve, StateReport, Status, StreamStatus } from '../machinepb/machine_pb';
import MqttContext from '../util/mqttContext'
import { ButtonGroup, Button, Typography, Slider, Box, Tabs, Tab } from '@mui/material';
import { useSessionStatus, useStateReport, useStreamStatus } from '../util/hooks';

export default function CollectDispense() {
    const noVials = 7
    const vials = new Array(noVials).fill(0).map((_, i) => noVials - i);

    const { client: c, messages } = useContext(MqttContext);
    const stateReport: StateReport | null = useStateReport();

    const [dropNumber, setDropNumber] = useState(3);

    const isAwake: boolean = !!stateReport &&
        stateReport?.getStatus() !== Status.SLEEPING &&
        stateReport?.getStatus() !== Status.E_STOP_ACTIVE;

    const collecting: boolean = !!stateReport && stateReport?.getCollectionRequest()?.getCompleted() === false;
    const collectingVial = collecting && stateReport?.getCollectionRequest()?.getVialNumber();

	// DROP VOLUMES
	const dropVolumeFromVial = (vial: number|undefined) => vial === 4 ? 12.0 : 14.0;

    const getAutoDispenseVolume = () => {
		return dropVolumeFromVial(stateReport?.getPipetteState()?.getVialHeld());
    }

    const requestCollection = (vial: number): void => {
        const volume = dropNumber * dropVolumeFromVial(vial);
        c?.publish(TOPIC_COLLECT, `${vial.toString()},${volume}`)
    }

    const keyDownHandler = (event: KeyboardEvent) => {
        const key = event.key;

        switch (key) {
            case ' ':
				c?.publish(TOPIC_DISPENSE, getAutoDispenseVolume().toString())
                break;
        }
    };

	const getDispensesRemaining = () => {
		return (stateReport?.getPipetteState()?.getVolumeTargetUl() ?? 0) / getAutoDispenseVolume();
	}

    useEffect(() => {
        window.addEventListener('keydown', keyDownHandler);
        return () => {
            window.removeEventListener('keydown', keyDownHandler);
        };
    }, [c, stateReport?.getPipetteState()?.getVialHeld()]);

	return (
		<>
            <Typography variant="h6">Collection & Dispense</Typography>
            <Slider
                value={dropNumber}
                onChange={(e, value) => typeof value === "number" ? setDropNumber(value): null}
                min={1}
                max={10} // Adjust the max value according to your requirements
                step={1}
                marks={true}
                valueLabelDisplay="auto"
                valueLabelFormat={(value) => `${value}`}
                aria-label="Collection drops"
                sx={{margin: 2, width: "50%"}}
            />
            <ButtonGroup variant="outlined" aria-label="outlined button group" sx={{margin: 2}}>
                {vials.map((vial) => 
                    <Button
                        key={vial}
                        disabled={!isAwake || collecting}
                        variant={collectingVial === vial ? "contained": "outlined"}
                        onClick={() => requestCollection(vial)}
                    >
                        {vial}
                    </Button>
                )}
            </ButtonGroup>
			<Typography variant="body1">Dispenses remaining: {getDispensesRemaining()}</Typography>
			<Typography variant="body1">Auto-Dispense Volume: {getAutoDispenseVolume()}µl</Typography>
            <Button disabled={!isAwake || collecting || stateReport?.getPipetteState()?.getSpent()} onClick={() => c?.publish(TOPIC_DISPENSE, getAutoDispenseVolume().toString())} sx={{"margin": 2}}>Auto-Dispense</Button>
		</>
		);
	}