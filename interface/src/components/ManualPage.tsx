import { useContext, useEffect, useRef, useState } from "react";
import {
	Button,
	ButtonGroup,
	Card,
	CardContent,
	Chip,
	Slider,
	Stack,
	TextField,
	Typography,
} from "@mui/material";
import { Mode, Status } from "../machinepb/machine";
import {
	TOPIC_MANUAL_PIPETTE_POSITION_SET,
	TOPIC_MANUAL_PIPETTE_SPEED_SET,
	TOPIC_MANUAL_PITCH_POSITION_SET,
	TOPIC_MANUAL_PITCH_SPEED_SET,
	TOPIC_MANUAL_RING_POSITION_SET,
	TOPIC_MANUAL_RING_SPEED_SET,
	TOPIC_MANUAL_YAW_POSITION_SET,
	TOPIC_MANUAL_YAW_SPEED_SET,
	TOPIC_MANUAL_Z_POSITION_SET,
	TOPIC_MANUAL_Z_SPEED_SET,
	TOPIC_STATE_REPORT_REQUEST,
	TOPIC_TOGGLE_MANUAL,
	TOPIC_WAKE,
} from "../topics_firmware/topics_firmware";
import { useStateReport } from "../util/hooks";
import MqttContext from "../util/mqttContext";

// Matches setMaxSpeed() in firmware initSteppers (steps/sec)
const MAX_SPEED = 1250;

interface Axis {
	name: string;
	unit: string;
	// Mirrors UnitStepper min/max units in firmware app/state.cpp
	min: number;
	max: number;
	speedTopic: string;
	positionTopic: string;
}

const AXES: Axis[] = [
	{
		name: "Ring",
		unit: "°",
		min: 4.8,
		max: 195,
		speedTopic: TOPIC_MANUAL_RING_SPEED_SET,
		positionTopic: TOPIC_MANUAL_RING_POSITION_SET,
	},
	{
		name: "Yaw",
		unit: "°",
		min: -21.3,
		max: 198,
		speedTopic: TOPIC_MANUAL_YAW_SPEED_SET,
		positionTopic: TOPIC_MANUAL_YAW_POSITION_SET,
	},
	{
		name: "Pitch",
		unit: "°",
		min: -2.5,
		max: 90,
		speedTopic: TOPIC_MANUAL_PITCH_SPEED_SET,
		positionTopic: TOPIC_MANUAL_PITCH_POSITION_SET,
	},
	{
		name: "Z",
		unit: "mm",
		min: 1,
		max: 73,
		speedTopic: TOPIC_MANUAL_Z_SPEED_SET,
		positionTopic: TOPIC_MANUAL_Z_POSITION_SET,
	},
	{
		name: "Pipette",
		unit: "µl",
		min: 0,
		max: 1000,
		speedTopic: TOPIC_MANUAL_PIPETTE_SPEED_SET,
		positionTopic: TOPIC_MANUAL_PIPETTE_POSITION_SET,
	},
];

function AxisControl({
	axis,
	enabled,
	jogSpeed,
}: {
	axis: Axis;
	enabled: boolean;
	jogSpeed: number;
}) {
	const { client: c } = useContext(MqttContext);
	const [speed, setSpeed] = useState(0);
	const [position, setPosition] = useState((axis.min + axis.max) / 2);
	const jogging = useRef(false);

	const publishSpeed = (s: number) => c?.publish(axis.speedTopic, s.toString());

	const startJog = (dir: number) => {
		if (!enabled) return;
		jogging.current = true;
		publishSpeed(dir * jogSpeed);
	};
	const stopJog = () => {
		if (!jogging.current) return;
		jogging.current = false;
		publishSpeed(0);
	};

	// Make sure a jog always ends, even if the release happens off the button
	// (or the button is disabled mid-jog) or the window loses focus.
	useEffect(() => {
		window.addEventListener("pointerup", stopJog);
		window.addEventListener("blur", stopJog);
		return () => {
			window.removeEventListener("pointerup", stopJog);
			window.removeEventListener("blur", stopJog);
		};
	});

	const jogButton = (dir: number, label: string) => (
		<Button
			disabled={!enabled}
			onPointerDown={() => startJog(dir)}
			onPointerUp={stopJog}
			onPointerLeave={stopJog}
			onPointerCancel={stopJog}
			sx={{ touchAction: "none", userSelect: "none" }}
		>
			{label}
		</Button>
	);

	const positionValid = position >= axis.min && position <= axis.max;

	return (
		<Card variant="outlined" sx={{ width: 380 }}>
			<CardContent>
				<Typography variant="h6">
					{axis.name}{" "}
					<Typography component="span" variant="caption">
						({axis.min}–{axis.max} {axis.unit})
					</Typography>
				</Typography>

				<Typography variant="subtitle2" sx={{ mt: 1 }}>
					Jog (hold)
				</Typography>
				<ButtonGroup size="small" variant="contained" sx={{ my: 1 }}>
					{jogButton(-1, "◀ −")}
					{jogButton(1, "+ ▶")}
				</ButtonGroup>

				<Typography variant="subtitle2">Speed (steps/s)</Typography>
				<Slider
					disabled={!enabled}
					value={speed}
					min={-MAX_SPEED}
					max={MAX_SPEED}
					step={25}
					valueLabelDisplay="auto"
					marks={[{ value: 0, label: "0" }]}
					onChange={(_, v) => setSpeed(v as number)}
				/>
				<ButtonGroup size="small" variant="outlined">
					<Button disabled={!enabled} onClick={() => publishSpeed(speed)}>
						Set speed
					</Button>
					<Button
						disabled={!enabled}
						color="error"
						onClick={() => {
							setSpeed(0);
							publishSpeed(0);
						}}
					>
						Stop
					</Button>
				</ButtonGroup>

				<Typography variant="subtitle2" sx={{ mt: 2 }}>
					Position ({axis.unit})
				</Typography>
				<Slider
					disabled={!enabled}
					value={position}
					min={axis.min}
					max={axis.max}
					step={0.1}
					valueLabelDisplay="auto"
					onChange={(_, v) => setPosition(v as number)}
				/>
				<Stack direction="row" spacing={1} alignItems="center">
					<TextField
						size="small"
						type="number"
						disabled={!enabled}
						error={!positionValid}
						value={position}
						inputProps={{ min: axis.min, max: axis.max, step: 0.1 }}
						onChange={(e) => setPosition(parseFloat(e.target.value))}
						sx={{ width: 120 }}
					/>
					<Button
						size="small"
						variant="outlined"
						disabled={!enabled || !positionValid}
						onClick={() => c?.publish(axis.positionTopic, position.toString())}
					>
						Go
					</Button>
				</Stack>
			</CardContent>
		</Card>
	);
}

export default function ManualPage() {
	const { client: c } = useContext(MqttContext);
	const stateReport = useStateReport();
	const [jogSpeed, setJogSpeed] = useState(300);

	useEffect(() => {
		if (!c || !c.connected) return;
		c.publish(TOPIC_STATE_REPORT_REQUEST, "");
	}, [c?.connected]);

	const isManual = stateReport?.mode === Mode.MANUAL;
	const isAwake =
		!!stateReport &&
		stateReport.status !== Status.SLEEPING &&
		stateReport.status !== Status.E_STOP_ACTIVE;

	const stopAll = () => AXES.forEach((a) => c?.publish(a.speedTopic, "0"));

	return (
		<div style={{ padding: "1rem" }}>
			<Typography variant="h5">Manual Control</Typography>

			<Stack direction="row" spacing={2} alignItems="center" sx={{ my: 2 }}>
				<Chip
					label={`Mode: ${stateReport ? Mode[stateReport.mode] : "UNKNOWN"}`}
					color={isManual ? "warning" : "default"}
				/>
				<Chip
					label={`Status: ${stateReport ? Status[stateReport.status] : "UNKNOWN"}`}
					color={isAwake ? "success" : "default"}
				/>
				<Button
					variant="outlined"
					disabled={isAwake}
					onClick={() => c?.publish(TOPIC_WAKE, "")}
				>
					Wake
				</Button>
				<Button
					variant="contained"
					color={isManual ? "secondary" : "warning"}
					disabled={!isAwake}
					onClick={() => c?.publish(TOPIC_TOGGLE_MANUAL, "")}
				>
					{isManual ? "Exit manual mode" : "Enter manual mode"}
				</Button>
				<Button
					variant="contained"
					color="error"
					disabled={!isManual}
					onClick={stopAll}
				>
					Stop all
				</Button>
			</Stack>

			<Stack direction="row" spacing={2} alignItems="center" sx={{ maxWidth: 500, mb: 2 }}>
				<Typography variant="subtitle2" noWrap>
					Jog speed (steps/s)
				</Typography>
				<Slider
					value={jogSpeed}
					min={25}
					max={MAX_SPEED}
					step={25}
					valueLabelDisplay="auto"
					onChange={(_, v) => setJogSpeed(v as number)}
				/>
			</Stack>

			{!isAwake && (
				<Typography variant="body2" color="error" sx={{ mb: 2 }}>
					Machine is asleep: the firmware ignores everything except wake and
					state report requests, so wake it before toggling manual mode.
				</Typography>
			)}

			{!isManual && (
				<Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
					Controls are disabled until the state report shows MANUAL mode (the
					firmware ignores manual commands otherwise).
				</Typography>
			)}

			<Stack direction="row" flexWrap="wrap" gap={2}>
				{AXES.map((a) => (
					<AxisControl key={a.name} axis={a} enabled={isManual} jogSpeed={jogSpeed} />
				))}
			</Stack>
		</div>
	);
}
