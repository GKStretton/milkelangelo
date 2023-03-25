import React, { useState, useContext, useEffect } from 'react';
import { TOPIC_SET_IK_Z, TOPIC_STATE_REPORT_JSON, TOPIC_STATE_REPORT_REQUEST } from '../util/topics'
import MqttContext from '../util/mqttContext'
import { Grid, ButtonGroup, Button } from '@mui/material';

function VerticalSlider() {
  const [value, setValue] = useState(42);

	const { client: c, messages } = useContext(MqttContext);
	const stateReportRaw = messages[TOPIC_STATE_REPORT_JSON];
  const stateReport = stateReportRaw && JSON.parse(stateReportRaw);
  console.log(stateReport);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = Number(event.target.value);
    setValue(newValue);
  };
  useEffect(() => {
    c?.publish(TOPIC_SET_IK_Z, value.toString());
  }, [value, c?.connected, stateReport?.status === "SLEEPING"]);

  return (
    <Grid
      container
      justifyContent="flex-start"
      alignItems="center"
      spacing={2}
    >
      <Grid item>
        <input
          type="range"
          min={35}
          max={50}
          value={value}
          step={1}
          onChange={handleChange}
          style={{
            height: '200px',
            width: '200px',
            transform: 'rotate(270deg)', /* rotate input 270 degrees */
          }}
        />
      </Grid>
      <Grid item>
        <ButtonGroup variant="contained" aria-label="text button group">
          <Button onClick={() => setValue((v) => v+1)}>+</Button>
          <Button onClick={() => setValue((v) => v-1)}>-</Button>
        </ButtonGroup>
      </Grid>
      <Grid item>
        {value}mm / <strong>{stateReport?.movement_details?.target_z_ik}mm</strong>
      </Grid>
    </Grid>
  );
}

export default VerticalSlider;