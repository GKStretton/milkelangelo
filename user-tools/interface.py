 
#!
#!
#!
#!
#! Make "setcrop" a single interface. use keys to switch between modes
#! crop mode, ik positioning mode, etc.
#!
#!
#!

import cv2

WINDOW = "window"
mouse_x = 0
mouse_y = 0

def mouse_callback(event, x, y, flags, param):
	global mouse_x, mouse_y
	if flags & cv2.EVENT_FLAG_LBUTTON:
		mouse_x = x
		mouse_y = y

# tcp is default
# top-cam-crop not working with ffmpeg pipelines for some reason. And gstreamer isn't installed
vcap = cv2.VideoCapture("rtsp://localhost:8554/top-cam")#, cv2.CAP_FFMPEG)

# Set up window settings
cv2.namedWindow(WINDOW)
cv2.setMouseCallback(WINDOW, mouse_callback)

while 1:
	ret, frame = vcap.read()
	if ret == False:
		print("Frame empty")
		break

	cv2.circle(frame,(mouse_x,mouse_y),10,(0,0,255),2, cv2.LINE_AA)

	cv2.imshow(WINDOW, frame)
	if cv2.waitKey(1) == 27:
		break