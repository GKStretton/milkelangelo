import React, { useState, useEffect, useContext } from 'react';
import { WebRTCReceiver } from '../util/WebRTCReceiver';
import { Button, Box, Typography } from '@mui/material';

interface VideoPlayerProps {
  url: string;
  name: string;
  handleClick?: (e: React.MouseEvent<HTMLVideoElement>) => void;
  renderOverlay?: (videoDimensions: { width: number; height: number }) => React.ReactNode;
  show: boolean;
}

const VideoPlayer: React.FC<VideoPlayerProps> = ({ url, name, handleClick, renderOverlay, show }: VideoPlayerProps) => {
  const [videoDimensions, setVideoDimensions] = useState<{ width: number; height: number }>({ width: 0, height: 0 });

  useEffect(() => {
    const receiver = new WebRTCReceiver(url, name);
  }, []);

  useEffect(() => {
    const videoElement = document.getElementById(name) as HTMLVideoElement;

    const updateVideoDimensions = () => {
      if (videoElement) {
        setVideoDimensions({ width: videoElement.clientWidth, height: videoElement.clientHeight });
      }
    };

    const resizeObserver = new ResizeObserver(() => {
      updateVideoDimensions();
    });

    if (videoElement) {
      updateVideoDimensions();
      resizeObserver.observe(videoElement);
    }

    return () => {
      if (videoElement) {
        resizeObserver.unobserve(videoElement);
      }
    };
  }, [name]);


  return (
    <Box display="flex" flexDirection="column" alignItems="center">
      <div style={{ position: 'relative' }}>
        <video
          id={name}
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
        {renderOverlay && renderOverlay(videoDimensions)}
      </div>
      <Box display="flex" flexDirection="row" alignItems="center">
        {/* <Typography>{JSON.stringify(videoDimensions)}</Typography> */}
      </Box>
    </Box>
  );
};

export default VideoPlayer;
