#!/bin/bash

# must be called after camera is being read.
/scripts/configure-top.sh &

gst-launch-1.0 v4l2src device=/dev/top-cam ! image/jpeg,width=1920,height=1080,framerate=30/1,format=MJPG ! rtspclientsink location=rtsp://localhost:8554/top-cam protocols=tcp