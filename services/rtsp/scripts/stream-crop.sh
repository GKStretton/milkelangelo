#!/bin/bash

# "top-cam" or "front-cam"
CAMERA=$1
CONFIG_FILE="/config/crop_${CAMERA}"
HOST="localhost"
SCRIPT_PATH="/scripts/getcrop.sh"

echo $CONFIG_FILE

top=$($SCRIPT_PATH $CONFIG_FILE top_rel)
bottom=$($SCRIPT_PATH $CONFIG_FILE bottom_rel)
right=$($SCRIPT_PATH $CONFIG_FILE right_rel)
left=$($SCRIPT_PATH $CONFIG_FILE left_rel)
echo $top $bottom $right $left

gst-launch-1.0 \
  rtspsrc protocols=tcp location="rtsp://$HOST:8554/$CAMERA" latency=0 ! \
  queue ! \
  rtph264depay ! \
  avdec_h264 ! \
  videocrop top=$top bottom=$bottom right=$right left=$left ! \
  x264enc bitrate=1500 speed-preset=ultrafast tune=zerolatency key-int-max=30 option-string="keyint_min=0" ! \
  rtspclientsink location=rtsp://$HOST:8554/${CAMERA}-crop protocols=tcp

  # tune=zerolatency has been removed from x264enc


