import React, { useContext } from 'react'
import './CleaningPage.css';
import MqttContext from '../util/mqttContext';
import { TOPIC_FLUID, TOPIC_MAINTENANCE, TOPIC_RINSE, TOPIC_SHUTDOWN, TOPIC_WAKE } from '../topics_firmware/topics_firmware';
import { FluidType } from '../machinepb/machine';
import { Button, Typography } from '@mui/material';

function CleaningPage() {
  const { client: c } = useContext(MqttContext);
  return (
    <div style={{padding: "1rem"}}>
      <br/>
      <Typography variant="h5">Cleaning Page</Typography>
      <ol>
        <li><input type="checkbox"/>
          <Button variant="contained" onClick={() => {
            c?.publish(TOPIC_WAKE, "")
            setTimeout(()=>{
              c?.publish(TOPIC_MAINTENANCE, "")
            }, 5000)
          }}>Turn on for maintenance</Button>
        </li>
        <li><input type="checkbox"/>Fill altar with 250ml cleaning solution</li>
        <li><input type="checkbox"/>Fill water bucket with 300ml cleaning solution</li>
        <li><input type="checkbox"/>
          <Button variant="contained" onClick={() => {
            c?.publish(TOPIC_FLUID, `${FluidType.FLUID_MILK},200,false`)
            setTimeout(()=>{
              c?.publish(TOPIC_FLUID, `${FluidType.FLUID_DRAIN},200,false`)
            }, 30000)
          }}>Run Cleaning Cycle</Button>
        </li>
        <li><input type="checkbox"/>To ensure pipette is unblocked, {' '}
          <Button variant="contained" onClick={() => {
            c?.publish(TOPIC_RINSE, "")
          }}>Rinse</Button>
          {' . Then, '}
          <Button variant="contained" color="secondary" onClick={() => {
            c?.publish(TOPIC_MAINTENANCE, "")
          }}>Free glass</Button>
        </li>
        <li><input type="checkbox"/>Change the rinse glass water</li>
        <li><input type="checkbox"/>Ensure vials are topped up with correct fluids</li>
        <li><input type="checkbox"/>Wait for cycle to complete (bowl to drain)</li>
        <li><input type="checkbox"/>Wipe bowl clean.{' '}
          <Button variant="contained" onClick={() => {
            c?.publish(TOPIC_RINSE, "")
          }}>Rinse</Button>
          {' '}To give room for wiping.
        </li>
        <li><input type="checkbox"/>Empty bucket and rinse/clean it</li>
        <li><input type="checkbox"/>Replace bucket, <strong>with tube in it</strong></li>
        <li><input type="checkbox"/>
          <Button variant="contained" onClick={() => {
            c?.publish(TOPIC_SHUTDOWN, "")
          }}>Turn off</Button>
        </li>
      </ol>
      <br/>
      <hr/>
      <br/>
      <Button variant="outlined" color="warning" onClick={() => {
        c?.publish(TOPIC_FLUID, `${FluidType.FLUID_DRAIN},50,false`)
      }}>Drain bowl more</Button>
    </div>
  )
}

export default CleaningPage;