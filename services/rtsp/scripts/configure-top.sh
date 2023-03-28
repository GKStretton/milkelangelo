#!/bin/bash

path=/dev/top-cam

echo "Starting $path configuration..."

v4l2-ctl --device $path -v width=1920,height=1080,pixelformat=MJPG --set-parm 30
v4l2-ctl --device $path --get-fmt-video --get-parm

v4l2-ctl --device $path --set-ctrl \
brightness=128,\
contrast=128,\
saturation=160,\
sharpness=128,\
backlight_compensation=0\

v4l2-ctl --device $path --set-ctrl \
focus_auto=0,\
exposure_auto=1,\
white_balance_temperature_auto=0,\
power_line_frequency=1

# have to set focus twice because it doesn't engage if you set it to the
# value as it already is. Seems to just stay at a silent 0 even though the value
# changes
v4l2-ctl --device $path --set-ctrl \
focus_absolute=29,\
white_balance_temperature=5000

v4l2-ctl --device $path --set-ctrl \
gain=1,\
zoom_absolute=120,\
exposure_absolute=70

# wait so focus set happens while camera is up (streamcam bug?)
sleep 1.5
v4l2-ctl --device $path --set-ctrl focus_absolute=30

# Fix for red camera effect
v4l2-ctl --device $path --set-ctrl gain=255
v4l2-ctl --device $path --set-ctrl gain=0

v4l2-ctl --device $path --list-ctrls-menus

echo "$path configuration done."