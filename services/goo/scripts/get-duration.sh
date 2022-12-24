#!/bin/bash
# returns duration in seconds of video file to stdout

ffprobe -i $1 -show_entries format=duration -v quiet -of csv="p=0"