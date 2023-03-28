#!/bin/bash

# "top-cam" or "front-cam"
CAMERA=$1
CONFIG_FILE="/mnt/md0/light-stores/kv/crop_${CAMERA}"
HOST="DEPTH"
SCRIPT_PATH="./getcrop.sh"

echo $CONFIG_FILE

top=$($SCRIPT_PATH $CONFIG_FILE top_rel)
bottom=$($SCRIPT_PATH $CONFIG_FILE bottom_rel)
right=$($SCRIPT_PATH $CONFIG_FILE right_rel)
left=$($SCRIPT_PATH $CONFIG_FILE left_rel)
echo $top $bottom $right $left

gst-launch-1.0 -v \
  rtspsrc protocols=tcp location="rtsp://$HOST:8554/$CAMERA" latency=0 ! \
  queue ! \
  rtph264depay ! \
  avdec_h264 ! \
  videocrop top=$top bottom=$bottom right=$right left=$left ! \
  x264enc bitrate=10000 speed-preset=ultrafast tune=zerolatency key-int-max=60 option-string="keyint_min=0" ! \
  rtspclientsink location=rtsp://$HOST:8554/${CAMERA}-crop protocols=tcp


