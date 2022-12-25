#!/bin/bash
# captures the stream $1 to the file $2
# NOTE: this saves the bitstream as-is. Compression / quality
# reduction I think requires re-encoding, so copy couldn't be used

ffmpeg -i rtsp://$1 -vcodec copy -r 30 -y $2