#!/bin/bash
# Generate post for all dslr images in session $1
CONTENT_DIR=/mnt/md0/light-stores/session_content/
python3 user-tools/auto_image_post.py -input ${CONTENT_DIR}$1/dslr/raw -output ${CONTENT_DIR}/$1/dslr/post