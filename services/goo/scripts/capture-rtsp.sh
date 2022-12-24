#!/bin/bash
# captures the stream $1 to the file $2

ffmpeg -i rtsp://$1 -b 100k -vcodec copy -r 30 -y $2