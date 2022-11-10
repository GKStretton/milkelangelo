#!/bin/bash

top=$(/scripts/getcrop.sh /config/crop.yml top_rel)
bottom=$(/scripts/getcrop.sh /config/crop.yml bottom_rel)
right=$(/scripts/getcrop.sh /config/crop.yml right_rel)
left=$(/scripts/getcrop.sh /config/crop.yml left_rel)
echo $top $bottom $right $left

gst-launch-1.0 -v rtspsrc protocols=tcp location="rtsp://localhost:8554/top-cam" latency=0 ! queue ! decodebin !  videocrop top=$top bottom=$bottom right=$right left=$left ! videoconvert ! rtspclientsink location=rtsp://localhost:8554/top-cam-crop protocols=tcp
