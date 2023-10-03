import { useState, useContext, useEffect } from "react";
import {
  TOPIC_SHUTDOWN,
  TOPIC_SLEEP,
  TOPIC_WAKE,
  TOPIC_SET_VALVE,
  TOPIC_DISPENSE,
  TOPIC_FLUID,
  TOPIC_COLLECT,
  TOPIC_SET_IK_Z,
  TOPIC_MARK_SAFE_TO_CALIBRATE,
  TOPIC_GOTO_NODE,
  TOPIC_TOGGLE_MANUAL,
  TOPIC_SET_COVER_OPEN,
  TOPIC_SET_COVER_CLOSE,
  TOPIC_RINSE,
  TOPIC_MAINTENANCE,
} from "../topics_firmware/topics_firmware";
import {
  TOPIC_SESSION_BEGIN,
  TOPIC_SESSION_END,
  TOPIC_SESSION_PAUSE,
  TOPIC_SESSION_RESUME,
  TOPIC_STREAM_END,
  TOPIC_STREAM_START,
} from "../topics_backend/topics_backend";
import { SessionStatus, SolenoidValve, StateReport, Status, StreamStatus, FluidType, Node } from "../machinepb/machine";
import MqttContext from "../util/mqttContext";
import {
  ButtonGroup,
  Button,
  Typography,
  Slider,
  Box,
  Tabs,
  Tab,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from "@mui/material";
import { useSessionStatus, useStateReport, useStreamStatus } from "../util/hooks";
import CollectDispense from "./CollectDispense";
import Profiles from "./Profiles";
import { useError } from "./ErrorManager";

export default function ControlGroup() {
  const error = useError();
  const [tabValue, setTabValue] = useState(0);
  const handleChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  const { client: c, messages } = useContext(MqttContext);
  const stateReport: StateReport | null = useStateReport();
  const sessionStatus: SessionStatus | null = useSessionStatus();
  const streamStatus: StreamStatus | null = useStreamStatus();

  const [dispenseVolume, setDispenseVolume] = useState(10.0);
  const [collectionVolume, setCollectionVolume] = useState(30.0);
  const [bulkFluidRequestVolume, setBulkFluidRequestVolume] = useState(200.0);
  const [zLevel, setZLevel] = useState(48);

  useEffect(() => {
    c?.publish(TOPIC_SET_IK_Z, zLevel.toString());
  }, [zLevel]);

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

  const isAwake: boolean =
    !!stateReport && stateReport?.status !== Status.SLEEPING && stateReport?.status !== Status.E_STOP_ACTIVE;

  const noVials = 6;
  const vials = new Array(noVials).fill(0).map((_, i) => noVials - i);

  const collecting: boolean = !!stateReport && stateReport?.collectionRequest?.completed === false;
  const collectingVial = collecting && stateReport?.collectionRequest?.vialNumber;

  const bulkRequests = [
    { id: 1, name: "Milk", fluid_type: FluidType.FLUID_MILK, open_drain: false },
    { id: 4, name: "Drain", fluid_type: FluidType.FLUID_DRAIN, open_drain: false },
    { id: 3, name: "Rinse", fluid_type: FluidType.FLUID_WATER, open_drain: true },
    { id: 2, name: "Water", fluid_type: FluidType.FLUID_WATER, open_drain: false },
  ];

  const valves = [
    { id: SolenoidValve.VALVE_DRAIN, name: "DRAIN" },
    { id: SolenoidValve.VALVE_WATER, name: "WATER" },
    { id: SolenoidValve.VALVE_MILK, name: "MILK" },
    { id: SolenoidValve.VALVE_AIR, name: "AIR" },
  ];

  const keyDownHandler = (event: KeyboardEvent) => {
    const key = event.key;

    switch (key) {
      case "w":
        c?.publish(TOPIC_WAKE, "");
        break;
      case "s":
        c?.publish(TOPIC_SHUTDOWN, "");
        break;
      case "k":
        c?.publish(TOPIC_SLEEP, "");
        break;
    }
  };

  useEffect(() => {
    window.addEventListener("keydown", keyDownHandler);
    return () => {
      window.removeEventListener("keydown", keyDownHandler);
    };
  }, [c, dispenseVolume, collectionVolume]);

  const [selectedNode, setSelectedNode] = useState<number>(Node.UNDEFINED);

  return (
    <>
      <Tabs value={tabValue} onChange={handleChange}>
        <Tab label="Core" />
        <Tab label="Overrides" />
        <Tab label="Profiles" />
        <Tab label="Sessions" />
        <Tab label="Socials" />
      </Tabs>

      {tabValue === 0 && (
        <>
          <Typography variant="h6">Basic</Typography>
          <ButtonGroup size="small" variant="outlined" aria-label="outlined button group" sx={{ margin: 1 }}>
            <Button disabled={!streamStatus || streamStatus.live} onClick={() => c?.publish(TOPIC_STREAM_START, "")}>
              Start Stream
            </Button>
          </ButtonGroup>
          <ButtonGroup size="small" variant="outlined" aria-label="outlined button group" sx={{ margin: 1 }}>
            <Button
              disabled={!sessionStatus || !sessionStatus.complete || sessionStatus.paused}
              color="warning"
              onClick={() => c?.publish(TOPIC_SESSION_BEGIN, "")}
            >
              {" "}
              Begin dev Session
            </Button>
            <Button
              disabled={!sessionStatus || !sessionStatus.complete || sessionStatus.paused}
              onClick={() => c?.publish(TOPIC_SESSION_BEGIN, "PRODUCTION")}
            >
              Begin Prod. Session
            </Button>
            <Button disabled={isAwake} onClick={() => c?.publish(TOPIC_WAKE, "")}>
              Wake
            </Button>
            <Button
              disabled={!isAwake}
              onClick={() =>
                c?.publish(TOPIC_FLUID, `${FluidType.FLUID_MILK},${bulkFluidRequestVolume},${false}`)
              }
            >
              Milk
            </Button>
          </ButtonGroup>
          <ButtonGroup size="small" variant="outlined" aria-label="outlined button group" sx={{ margin: 1 }}>
            <Button
              disabled={!sessionStatus || sessionStatus.complete || sessionStatus.paused}
              onClick={() => c?.publish(TOPIC_SESSION_PAUSE, "")}
            >
              Pause Session
            </Button>
            <Button
              disabled={!sessionStatus || !sessionStatus.paused}
              onClick={() => c?.publish(TOPIC_SESSION_RESUME, "")}
            >
              Resume Session
            </Button>
          </ButtonGroup>
          <ButtonGroup size="small" variant="outlined" aria-label="outlined button group" sx={{ margin: 1 }}>
            <Button
              disabled={!isAwake}
              onClick={() =>
                c?.publish(TOPIC_FLUID, `${FluidType.FLUID_DRAIN},${bulkFluidRequestVolume},${false}`)
              }
            >
              Drain
            </Button>
            <Button
              disabled={!isAwake}
              onClick={() =>
                c?.publish(TOPIC_FLUID, `${FluidType.FLUID_WATER},${bulkFluidRequestVolume},${true}`)
              }
            >
              Rinse
            </Button>
            <Button disabled={!isAwake} onClick={() => c?.publish(TOPIC_SHUTDOWN, "")}>
              Shutdown
            </Button>
            <Button disabled={!isAwake} variant="contained" color="error" onClick={() => c?.publish(TOPIC_SLEEP, "")}>
              Kill
            </Button>
            <Button
              disabled={!sessionStatus || sessionStatus.complete}
              onClick={() => c?.publish(TOPIC_SESSION_END, "")}
            >
              End Session
            </Button>
            <Button disabled={!streamStatus || !streamStatus.live} onClick={() => c?.publish(TOPIC_STREAM_END, "")}>
              End Stream
            </Button>
          </ButtonGroup>

          <Typography variant="h6">Bulk Fluid</Typography>
          <Slider
            value={bulkFluidRequestVolume}
            onChange={(e, value) => (typeof value === "number" ? setBulkFluidRequestVolume(value) : null)}
            min={25}
            max={300} // Adjust the max value according to your requirements
            step={25}
            marks={marks_ml}
            valueLabelDisplay="auto"
            valueLabelFormat={(value) => `${value}ml`}
            aria-label="Bulk Fluid Request volume"
            sx={{ margin: 2, width: "50%" }}
          />
          <ButtonGroup variant="outlined" aria-label="outlined button group" sx={{ margin: 2 }}>
            {bulkRequests.map((request) => (
              <Button
                key={request.id}
                disabled={!isAwake}
                onClick={() =>
                  c?.publish(TOPIC_FLUID, `${request.fluid_type},${bulkFluidRequestVolume},${request.open_drain}`)
                }
              >
                {request.name}
              </Button>
            ))}
          </ButtonGroup>

          <CollectDispense />

          <Typography variant="h6">Z-Level</Typography>
          <Slider
            value={zLevel}
            onChange={(e, value) => (typeof value === "number" ? setZLevel(value) : null)}
            min={35}
            max={55} // Adjust the max value according to your requirements
            step={1}
            marks={zMarks}
            valueLabelDisplay="auto"
            valueLabelFormat={(value) => `${value}mm`}
            aria-label="Z-Level"
            sx={{ margin: 2, width: "50%" }}
          />
        </>
      )}

      {tabValue === 1 && (
        <>
          <Typography variant="h6">Open Valves</Typography>
          <Box display="flex" flexDirection="row" alignItems="center" justifyContent="center" sx={{ margin: 0.5 }}>
            <ButtonGroup size="small" sx={{ margin: 0.5 }}>
              {valves.map((valve) => (
                <Button
                  key={valve.id}
                  disabled={!isAwake}
                  onClick={() => c?.publish(TOPIC_SET_VALVE, `${valve.id},true`)}
                >
                  {valve.name} (O)
                </Button>
              ))}
            </ButtonGroup>
          </Box>
          <Typography variant="h6">Close Valves</Typography>
          <Box display="flex" flexDirection="row" alignItems="center" justifyContent="center" sx={{ margin: 0.5 }}>
            <ButtonGroup size="small" sx={{ margin: 0.5 }}>
              {valves.map((valve) => (
                <Button
                  key={valve.id}
                  disabled={!isAwake}
                  onClick={() => c?.publish(TOPIC_SET_VALVE, `${valve.id},false`)}
                >
                  {valve.name} (C)
                </Button>
              ))}
            </ButtonGroup>
          </Box>

          <Typography variant="h6">Collection</Typography>
          <Slider
            value={collectionVolume}
            onChange={(e, value) => (typeof value === "number" ? setCollectionVolume(value) : null)}
            min={1}
            max={100} // Adjust the max value according to your requirements
            step={1}
            marks={marks}
            valueLabelDisplay="auto"
            valueLabelFormat={(value) => `${value}µl`}
            aria-label="Dispense volume"
            sx={{ margin: 2, width: "50%" }}
          />
          <ButtonGroup variant="outlined" aria-label="outlined button group" sx={{ margin: 2 }}>
            {vials.map((vial) => (
              <Button
                key={vial}
                disabled={!isAwake || collecting}
                variant={collectingVial === vial ? "contained" : "outlined"}
                onClick={() => c?.publish(TOPIC_COLLECT, `${vial.toString()},${collectionVolume}`)}
              >
                {vial}
              </Button>
            ))}
          </ButtonGroup>

          <Typography variant="h6">Dispense</Typography>
          <Slider
            value={dispenseVolume}
            onChange={(e, value) => (typeof value === "number" ? setDispenseVolume(value) : null)}
            min={1}
            max={collectionVolume} // Adjust the max value according to your requirements
            step={1}
            marks={marks}
            valueLabelDisplay="auto"
            valueLabelFormat={(value) => `${value}µl`}
            aria-label="Dispense volume"
            sx={{ margin: 2, width: "50%" }}
          />

          <Button
            disabled={!isAwake || stateReport?.pipetteState?.spent}
            onClick={() => c?.publish(TOPIC_DISPENSE, dispenseVolume.toString())}
            sx={{ margin: 2 }}
          >
            Dispense
          </Button>

          <Typography variant="h6">Other</Typography>
          <Button
            disabled={!isAwake}
            variant="contained"
            onClick={() => c?.publish(TOPIC_MAINTENANCE, "")}
            sx={{ margin: 2 }}
          >
            Goto Maintenance Position
          </Button>
          <Button
            disabled={!isAwake}
            variant="contained"
            color="error"
            onClick={() => c?.publish(TOPIC_MARK_SAFE_TO_CALIBRATE, dispenseVolume.toString())}
            sx={{ margin: 2 }}
          >
            Mark Safe to Calibrate
          </Button>

          <FormControl margin="normal">
            <InputLabel id="node-select-label">Node</InputLabel>
            <Select
              labelId="node-select-label"
              label="Node"
              value={selectedNode}
              onChange={(event) => {
                const node = event.target.value as number;
                setSelectedNode(node);
                console.log(`selected node ${node}`);
              }}
            >
              {Object.entries(Node).map(([key, value]) => (
                <MenuItem value={value} key={key}>
                  {key}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          <Button
            disabled={!isAwake}
            variant="contained"
            color="error"
            onClick={() => c?.publish(TOPIC_GOTO_NODE, String(selectedNode))}
          >
            Go to node
          </Button>

          <Button
            disabled={!isAwake}
            variant="contained"
            color="error"
            onClick={() => c?.publish(TOPIC_TOGGLE_MANUAL, "")}
            sx={{ margin: 2 }}
          >
            Toggle Manual
          </Button>

          <Button
            disabled={!isAwake}
            variant="contained"
            color="error"
            onClick={() => c?.publish(TOPIC_SET_COVER_OPEN, "")}
            sx={{ margin: 2 }}
          >
            Open Cover
          </Button>

          <Button
            disabled={!isAwake}
            variant="contained"
            color="error"
            onClick={() => c?.publish(TOPIC_SET_COVER_CLOSE, "")}
            sx={{ margin: 2 }}
          >
            Close Cover
          </Button>

          <Button
            variant="contained"
            color="error"
            onClick={() => {
              error("test");
            }}
            sx={{ margin: 2 }}
          >
            Test Error
          </Button>

          <Button
            disabled={!isAwake}
            variant="contained"
            color="error"
            onClick={() => c?.publish(TOPIC_RINSE, "")}
            sx={{ margin: 2 }}
          >
            Rinse
          </Button>
        </>
      )}

      {tabValue === 2 && <Profiles />}

      {tabValue === 3 && (
        <>
          <Typography variant="h4"> after session </Typography>
          <code>cd ~/asol/software/asol-backend</code>
          <ol>
            <li>choose image</li>
            <code>./scripts/interactive/feh-dslr-selection-utility $session</code>
            <li>generate content plan</li>
            <code>./scripts/generation/generate-content-plan $session</code>
            <li>generate content</li>
            <code>./scripts/generation/generate-content $session</code>
            <li>publish</li>
            <code>./scripts/interactive/publish-session $session</code>
          </ol>
          <Typography variant="h4"> video generation modifications </Typography>
          <Typography> If changes to the pipeline required: </Typography>

        </>
      )}

      {tabValue === 4 && (
        <>
          <Typography>Not implemented: post to social media</Typography>

          <Button
            onClick={() =>
              c?.publish("asol/videos-generated", sessionStatus?.id === 0 ? "" : String(sessionStatus?.id))
            }
          >
            Signal videos generated (do upload)
          </Button>
        </>
      )}
    </>
  );
}
