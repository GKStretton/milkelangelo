#!/bin/bash
# Capture and save a DSLR photo to $1

gphoto2\
	--capture-image-and-download \
	--force-overwrite \
	--filename $1
