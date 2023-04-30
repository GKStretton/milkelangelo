#!/bin/bash
# generate all video content for session $1

# python3 user-tools/auto_video_post.py -y -n $1 -t SHORTFORM -d /home/greg/Downloads #-d /mnt/md0/light-stores
python3 user-tools/auto_video_post.py -y -n $1 -t CLEANING -d /home/greg/Downloads #-d /mnt/md0/light-stores
# python3 user-tools/auto_dslr_timelapse.py -n $1 -d /home/greg/Downloads #-d /mnt/md0/light-stores
