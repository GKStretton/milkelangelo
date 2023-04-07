import React, { useState, useContext } from 'react';
import MqttContext from '../util/mqttContext'
import { TOPIC_GOTO_XY } from '../topics_firmware/topics_firmware';
import VideoPlayer from './VideoPlayer';

const TopCamPlayer = () => {
  const [circlePos, setCirclePos] = useState<{ x: number; y: number }>({ x: 0, y: 0 });
	const { client: c, messages } = useContext(MqttContext);

  const handleClick = (e: React.MouseEvent<HTMLVideoElement>) => {
    const rect = e.currentTarget.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;

    const normalizedX = (x / rect.width) * 2 - 1;
    const normalizedY = -((y / rect.height) * 2 - 1);

    console.log(`Clicked at normalized coordinates: (${normalizedX}, ${normalizedY})`);

    setCirclePos({ x: normalizedX, y: normalizedY });

    c?.publish(TOPIC_GOTO_XY, `${normalizedX},${normalizedY}`)
  };

  const renderOverlay = (videoDimensions: { width: number; height: number }) => (
    <>
    <div
      style={{
        position: 'absolute',
        border: '5px solid red',
        borderRadius: '50%',
        width: '20px',
        height: '20px',
        left: `${(circlePos.x + 1) * 0.5 * videoDimensions.width}px`,
        top: `${(-circlePos.y + 1) * 0.5 * videoDimensions.height}px`,
        transform: 'translate(-50%, -50%)',
      }}
    ></div>
    <img
      src="/mask_alpha.png"
      alt="alpha mask"
      style={{
        position: 'absolute',
        top: 0,
        left: 0,
        width: `${videoDimensions.width+1}px`,
        height: `${videoDimensions.height+1}px`,
      }}
    />
    </>
  );

  return (
    <>
      <div style={{ position: 'relative' }}>
        <VideoPlayer
        url="DEPTH:8889/top-cam-crop/"
          name="top"
          handleClick={handleClick}
          renderOverlay={renderOverlay}
          show={false}
        />
      </div>
    </>
  );
};

export default TopCamPlayer;
