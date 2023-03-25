import React, { useState, useContext } from 'react';
import { TOPIC_SET_IK_Z, TOPIC_STATE_REPORT_JSON, TOPIC_STATE_REPORT_REQUEST } from '../util/topics'
import MqttContext from '../util/mqttContext'

function VerticalSlider() {
  const [value, setValue] = useState(70);

	const { client: c, messages } = useContext(MqttContext);
	const stateReportRaw = messages[TOPIC_STATE_REPORT_JSON];
  const stateReport = stateReportRaw && JSON.parse(stateReportRaw);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = Number(event.target.value);
    setValue(newValue);
    c?.publish(TOPIC_SET_IK_Z, newValue.toString());
  };

  return (
    <div>
      <input
        type="range"
        min={60}
        max={80}
        value={value}
        step={1}
        onChange={handleChange}
        style={{
          height: '200px',
          transform: 'rotate(270deg)', /* rotate input 270 degrees */
        }}
      />
      <div style={{ textAlign: 'center' }}>
        {value}mm / {stateReport?.movement_details?.target_ik_z}
      </div>
    </div>
  );
}

export default VerticalSlider;