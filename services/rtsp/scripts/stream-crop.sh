#!/bin/bash

# "top-cam" or "front-cam"
CAMERA=$1
CONFIG_FILE="/config/crop_${CAMERA}"

top=$(/scripts/getcrop.sh $CONFIG_FILE top_rel)
bottom=$(/scripts/getcrop.sh $CONFIG_FILE bottom_rel)
right=$(/scripts/getcrop.sh $CONFIG_FILE right_rel)
left=$(/scripts/getcrop.sh $CONFIG_FILE left_rel)
echo $top $bottom $right $left

gst-launch-1.0 -v rtspsrc protocols=tcp location="rtsp://localhost:8554/$CAMERA" latency=0 ! queue ! decodebin ! videocrop top=$top bottom=$bottom right=$right left=$left ! rtspclientsink location=rtsp://localhost:8554/${CAMERA}-crop protocols=tcp
