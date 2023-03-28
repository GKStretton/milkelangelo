#!/bin/bash

# must be called after camera is being read.
/scripts/configure-front.sh &

gst-launch-1.0 v4l2src device=/dev/front-cam ! image/jpeg,width=1920,height=1080,framerate=30/1,format=MJPG ! jpegdec ! videoconvert ! x264enc bitrate=5000 speed-preset=ultrafast tune=zerolatency key-int-max=30 option-string="keyint_min=0" ! rtspclientsink location=rtsp://localhost:8554/front-cam protocols=tcp