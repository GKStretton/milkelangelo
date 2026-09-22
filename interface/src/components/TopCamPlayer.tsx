import React, { useState, useContext } from 'react';
import MqttContext from '../util/mqttContext'
import { TOPIC_GOTO_XY } from '../topics_firmware/topics_firmware';
import VideoPlayer from './VideoPlayer';

interface TopCamProps {
  url: string;
}

const TopCamPlayer = ({ url }: TopCamProps) => {
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
        width: '24px',
        height: '24px',
        left: `${(circlePos.x + 1) * 0.5 * videoDimensions.width}px`,
        top: `${(-circlePos.y + 1) * 0.5 * videoDimensions.height}px`,
        transform: 'translate(-50%, -50%)',
        pointerEvents: 'none',
      }}
    >
      <div
        style={{
          position: 'absolute',
          top: '50%',
          left: 0,
          width: '100%',
          height: '2px',
          background: 'red',
          transform: 'translateY(-50%)',
        }}
      ></div>
      <div
        style={{
          position: 'absolute',
          left: '50%',
          top: 0,
          height: '100%',
          width: '2px',
          background: 'red',
          transform: 'translateX(-50%)',
        }}
      ></div>
    </div>
    <img
      src="/mask_alpha.png"
      alt="alpha mask"
      style={{
        position: 'absolute',
        top: 0,
        left: 0,
        width: `${videoDimensions.width+1}px`,
        height: `${videoDimensions.height+1}px`,
        pointerEvents: 'none',
      }}
    />
    </>
  );

  return (
    <>
      <div style={{ position: 'relative' }}>
        <VideoPlayer
        url={`${url}/top-cam-crop/`}
          name="top"
          handleClick={handleClick}
          renderOverlay={renderOverlay}
          aspectRatio="1 / 1"
        />
      </div>
    </>
  );
};

export default TopCamPlayer;
