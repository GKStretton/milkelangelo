#!/bin/bash
# generate all stills for a session id

set -e

file=user-tools/auto_stills_editing.py
baseDir=/mnt/md0/light-stores
number=$1

if [ -z "$number" ]; then
	echo "Session not specified as argument. Exiting."
	exit 1
fi

if [ ! -d "$baseDir/session_content/$number" ]; then
	echo "Directory $baseDir/session_content/$number is not present. Exiting."
	exit 1
fi

python3 $file --base-dir $baseDir --session-number $number -t INTRO -f PORTRAIT
python3 $file --base-dir $baseDir --session-number $number -t INTRO -f LANDSCAPE
python3 $file --base-dir $baseDir --session-number $number -t OUTRO -f PORTRAIT
python3 $file --base-dir $baseDir --session-number $number -t OUTRO -f LANDSCAPE

echo "Files successfully output to $baseDir/session_content/$number/stills"