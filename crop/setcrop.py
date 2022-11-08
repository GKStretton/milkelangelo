# The crop tool lets you configure the crop on the top webcam. Its output is
# used by the cropper in rtsp.

import cv2
import yaml
import os

YML_FILE = "crop.yml"
WINDOW = "window"

# left click
x1 = 0
y1 = 0

# right click
x2 = 0
y2 = 0

if os.path.isfile(YML_FILE):
	with open(YML_FILE, 'r') as f:
		yml = yaml.load(f)
		x1 = yml['left_abs']
		y1 = yml['top_abs']
		x2 = yml['right_abs']
		y2 = yml['bottom_abs']
		print("loaded from file:", yml)

def mouse_callback(event, x, y, flags, param):
	global x1, y1, x2, y2
	if flags & cv2.EVENT_FLAG_LBUTTON and flags & cv2.EVENT_FLAG_SHIFTKEY:
		x2 = x
		y2 = y
	elif flags & cv2.EVENT_FLAG_LBUTTON:
		x1 = x
		y1 = y

# tcp is default
vcap = cv2.VideoCapture("rtsp://DEPTH:8554/top-cam")#, cv2.CAP_GSTREAMER)

# Set up window settings
cv2.namedWindow(WINDOW)
cv2.setMouseCallback(WINDOW, mouse_callback)

width = 0
height = 0

while 1:
	ret, frame = vcap.read()
	if ret == False:
		print("Frame empty")
		break

	cv2.rectangle(frame,(x1,y1), (x2, y1 + x2 - x1), (0,0,255),2, cv2.LINE_AA)

	width = frame.shape[1]
	height = frame.shape[0]

	cv2.imshow(WINDOW, frame)
	if cv2.waitKey(1) == 27:
		break

print(width, height)

result = {
	"left_abs": x1,
	"right_abs": x2,
	"top_abs": y1,
	"bottom_abs": y2,
	"left_rel": x1,
	"right_rel": width - x2,
	"top_rel": y1,
	"bottom_rel": height - y2,
}
print(result)
with open(YML_FILE, 'w') as f:
	yaml.dump(result, f)