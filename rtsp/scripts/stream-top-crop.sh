#!/bin/bash

top=$(/crop/getcrop.sh /crop/crop.yml top_rel)
bottom=$(/crop/getcrop.sh /crop/crop.yml bottom_rel)
right=$(/crop/getcrop.sh /crop/crop.yml right_rel)
left=$(/crop/getcrop.sh /crop/crop.yml left_rel)
echo $top $bottom $right $left
gst-launch-1.0 -v rtspsrc protocols=tcp location="rtsp://localhost:8554/top-cam" latency=0 ! queue ! decodebin ! videoconvert ! videocrop top=$top bottom=$bottom right=$right left=$left ! videoconvert ! rtspclientsink location=rtsp://localhost:8554/top-cam-crop protocols=tcp
