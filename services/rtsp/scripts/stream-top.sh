#!/bin/bash

# must be called after camera is being read.
/scripts/configure-top.sh &

gst-launch-1.0 v4l2src device=/dev/top-cam ! image/jpeg,width=1920,height=1080,framerate=60/1,format=MJPG ! jpegdec ! videoconvert ! x264enc bitrate=10000 speed-preset=ultrafast tune=zerolatency key-int-max=60 option-string="keyint_min=0" ! rtspclientsink location=rtsp://localhost:8554/top-cam protocols=tcp