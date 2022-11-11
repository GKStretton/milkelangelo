# The crop tool lets you configure the crop on the top webcam. Its output is
# used by the cropper in rtsp.

import cv2
import numpy as np
import yaml
import os
import sys
import paho.mqtt.publish as mqttpub
import paho.mqtt.subscribe as mqttsub

def overlay_image_alpha(img, img_overlay, x, y, alpha_mask):
    """Overlay `img_overlay` onto `img` at (x, y) and blend using `alpha_mask`.

    `alpha_mask` must have same HxW as `img_overlay` and values in range [0, 1].
    """
    # Image ranges
    y1, y2 = max(0, y), min(img.shape[0], y + img_overlay.shape[0])
    x1, x2 = max(0, x), min(img.shape[1], x + img_overlay.shape[1])

    # Overlay ranges
    y1o, y2o = max(0, -y), min(img_overlay.shape[0], img.shape[0] - y)
    x1o, x2o = max(0, -x), min(img_overlay.shape[1], img.shape[1] - x)

    # Exit if nothing to do
    if y1 >= y2 or x1 >= x2 or y1o >= y2o or x1o >= x2o:
        return

    # Blend overlay within the determined ranges
    img_crop = img[y1:y2, x1:x2]
    img_overlay_crop = img_overlay[y1o:y2o, x1o:x2o]
    alpha = alpha_mask[y1o:y2o, x1o:x2o, np.newaxis]
    alpha_inv = 1.0 - alpha

    img_crop[:] = alpha * img_overlay_crop + alpha_inv * img_crop

WINDOW = "window"

TOP_MASK = "./resources/static_img/top-mask.png"
GET_TOPIC = "crop-config/get"
GET_RESP_TOPIC = "crop-config/get-resp"
SET_TOPIC = "crop-config/set"
SET_RESP_TOPIC = "crop-config/set-resp"
CLIENT_ID="crop-setter-ui"

# left click
x1 = 0
y1 = 0

# right click
x2 = 0
y2 = 0

def load_config(yml):
	global x1, y1, x2, y2
	x1 = yml['left_abs']
	y1 = yml['top_abs']
	x2 = yml['right_abs']
	y2 = yml['bottom_abs']

	print("loaded values from yaml:", yml)

if len(sys.argv) > 1:
	YML_FILE = sys.argv[1]
	LOCAL = True
	print("local mode")
else:
	LOCAL = False
	print("remote mode")

# LOAD EXISTING CONFIG
if LOCAL:
	print("loading from", YML_FILE, "...")
	if os.path.isfile(YML_FILE):
		with open(YML_FILE, 'r') as f:
			yml = yaml.load(f)
			load_config(yml)
	else:
		print("file not found, proceeding with 0 values")
else:
	mqttpub.single(GET_TOPIC, payload="get", hostname="localhost", port=1883, client_id=CLIENT_ID)
	print("sent config request to", GET_TOPIC)
	print("waiting for response on", GET_RESP_TOPIC, "...")
	resp = mqttsub.simple(GET_RESP_TOPIC, hostname="localhost", port=1883, client_id=CLIENT_ID, keepalive=1)
	if resp.payload == "404":
		print("no config yet, using 0 values")
	else:
		print("got response", resp.payload)
		yml = yaml.load(resp.payload, Loader=yaml.FullLoader)
		load_config(yml)


def mouse_callback(event, x, y, flags, param):
	global x1, y1, x2, y2
	if flags & cv2.EVENT_FLAG_LBUTTON and flags & cv2.EVENT_FLAG_SHIFTKEY:
		x2 = x
		y2 = y
	elif flags & cv2.EVENT_FLAG_LBUTTON:
		x1 = x
		y1 = y

print("attempting rtsp conn...")
# tcp is default
vcap = cv2.VideoCapture("rtsp://localhost:8554/top-cam")#, cv2.CAP_GSTREAMER)
if not vcap.isOpened():
	print("Error loading webcam stream, aborting.")
	exit()

# Set up window settings
cv2.namedWindow(WINDOW)
cv2.setMouseCallback(WINDOW, mouse_callback)

width = 0
height = 0

ret, frame = vcap.read()
if ret == False:
	print("Frame empty")
	exit()

width = frame.shape[1]
height = frame.shape[0]

# load mask
mask = cv2.imread(TOP_MASK)
print(mask.shape)

print("starting loop")
while 1:
	ret, frame = vcap.read()
	if ret == False:
		print("Frame empty")
		break

	cv2.rectangle(frame,(x1,y1), (x2, y1 + x2 - x1), (0,0,255),2, cv2.LINE_AA)

	resized_mask = cv2.resize(mask, (x2-x1, x2-x1))
	res = frame.copy()
	# https://stackoverflow.com/a/45118011
	overlay_image_alpha(res, np.zeros((x2-x1,x2-x1,3)), x1, y1, 1 - resized_mask[:,:,0] / 255.0)


	cv2.imshow(WINDOW, res)
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

if LOCAL:
	with open(YML_FILE, 'w') as f:
		yaml.dump(result, f)
else:
	conf = yaml.dump(result)
	mqttpub.single(SET_TOPIC, payload=conf, hostname="localhost", port=1883, client_id=CLIENT_ID)
	print("sent conf to", SET_TOPIC)
	
	print("waiting for response...")
	resp = mqttsub.simple(SET_RESP_TOPIC, hostname="localhost", port=1883, client_id=CLIENT_ID, keepalive=1)
	print("got response", resp.payload.decode("utf-8"))