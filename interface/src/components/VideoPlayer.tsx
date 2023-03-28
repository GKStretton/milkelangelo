import React, { useState, useEffect, useContext } from 'react';
import MqttContext from '../util/mqttContext'
import { TOPIC_KV_GET, TOPIC_KV_GET_RESP } from '../util/topics';
import { WebRTCReceiver } from '../util/WebRTCReceiver';
import yaml from 'js-yaml';
import { Button, Box, Typography } from '@mui/material';

interface VideoPlayerProps {
  url: string;
  name: string;
}

interface CropConfig {
  top_rel: number;
  top_abs: number;
  bottom_rel: number;
  bottom_abs: number;
  right_rel: number;
  right_abs: number;
  left_rel: number;
  left_abs: number;
}

const VideoPlayer: React.FC<VideoPlayerProps> = ({ url, name }: VideoPlayerProps) => {
	const { client: c, messages } = useContext(MqttContext);
  const [circlePos, setCirclePos] = useState({ x: -1, y: -1 });
  const [doCrop, setDoCrop] = useState(true);

  const cropRaw = messages[`${TOPIC_KV_GET_RESP}crop_${name}-cam`];
  let crop: CropConfig | null = null;
  if (cropRaw) {
    crop = yaml.load(cropRaw) as CropConfig;
    console.log(`crop ${name}:`, crop);
  }

  useEffect(() => {
    const receiver = new WebRTCReceiver(url, name);
  }, []);

  useEffect(() => {
    c?.subscribe(`${TOPIC_KV_GET_RESP}crop_${name}-cam`);
    c?.publish(`${TOPIC_KV_GET}crop_${name}-cam`, "");
  }, [c?.connected])

  const handleClick = (e: React.MouseEvent<HTMLVideoElement>) => {
    const rect = e.currentTarget.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;

    console.log(`Clicked at coordinates: (${x}, ${y})`);

    setCirclePos({ x, y });
  };

  const width = 1920;
  const height = 1080;

  return (
    <Box display="flex" flexDirection="column" alignItems="flex-start">
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
            objectFit: 'cover',
            // clipPath: doCrop && crop ? `inset(${crop.top_rel / height * 100}% ${crop.right_rel / width * 100}% ${crop.bottom_rel / height * 100}% ${crop.left_rel / width * 100}%)` : 'none',
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
      <Box display="flex" flexDirection="row" alignItems="center">
        <Button onClick={()=>{setDoCrop(!doCrop)}}>{doCrop ? "Disable":"Enable"} Crop</Button>
      </Box>
    </Box>
  );
};

export default VideoPlayer;
