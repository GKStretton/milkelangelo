import React, { useState, useEffect } from 'react';
import { WebRTCReceiver } from '../util/WebRTCReceiver';

interface VideoPlayerProps {
  url: string;
}

const VideoPlayer: React.FC<VideoPlayerProps> = ({ url }: VideoPlayerProps) => {
  const [circlePos, setCirclePos] = useState({ x: -1, y: -1 });

  useEffect(() => {
    const receiver = new WebRTCReceiver(url);
  }, []);

  const handleClick = (e: React.MouseEvent<HTMLVideoElement>) => {
    const rect = e.currentTarget.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;

    console.log(`Clicked at coordinates: (${x}, ${y})`);

    setCirclePos({ x, y });
  };

  return (
    <div style={{ position: 'relative' }}>
      <video
        id="video"
        muted
        controls={false}
        onClick={handleClick}
        autoPlay
        playsInline
        style={{
          width: '100%',
          height: '100%',
          border: '1px solid black',
        }}
      ></video>
      {circlePos.x > -1 && (
        <div
          style={{
            position: 'absolute',
            border: '5px solid red',
            borderRadius: '50%',
            width: '20px',
            height: '20px',
            left: `${circlePos.x}px`,
            top: `${circlePos.y}px`,
            transform: 'translate(-50%, -50%)',
          }}
        ></div>
      )}
    </div>
  );
};

export default VideoPlayer;
