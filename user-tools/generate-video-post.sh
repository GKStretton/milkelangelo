#!/bin/bash
# generate all video content for session $1

# python3 user-tools/auto_video_post.py -y -n $1 -t CONTENT_TYPE_LONGFORM -d /mnt/md0/light-stores
python3 user-tools/auto_video_post.py -y -n $1 -t CONTENT_TYPE_SHORTFORM -d /mnt/md0/light-stores
python3 user-tools/auto_video_post.py -y -n $1 -t CONTENT_TYPE_CLEANING -d /mnt/md0/light-stores
python3 user-tools/auto_dslr_timelapse.py -n $1 -d /mnt/md0/light-stores
