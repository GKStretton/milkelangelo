import React, { useState } from 'react';

function VerticalSlider() {
  const [value, setValue] = useState(70);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = Number(event.target.value);
    setValue(newValue);
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
        {value}mm
      </div>
    </div>
  );
}

export default VerticalSlider;