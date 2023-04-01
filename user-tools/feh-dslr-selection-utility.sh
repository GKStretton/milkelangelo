#!/bin/bash
# script for selecting the final dslr image via feh for session $1
# key 1 will do the selection inside feh. arrow keys
# to navigate. q to quit. stills will be generated on quit

DIR=$(pwd)
DSLR_DIR=/mnt/md0/light-stores/session_content/$1/dslr/post
cd $DSLR_DIR
feh --action1 "; $DIR/user-tools/select-dslr-image.sh %F" --scale-down --draw-filename .

cd $DIR
./user-tools/generate-all-stills.sh $1