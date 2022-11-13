import cv2
import time
import math
from pycommon.window import Window
from pycommon.config_manager_client import read_remote_crop_config
import pycommon.image as image
import numpy as np
import pycommon.mqtt_client as mc

TOP_MASK = "../resources/static_img/top-mask.png"
# how many microlitres to dispense at a time
DISPENSE_uL = 10

helps = {
	"h": "Print this help text",
	"c": "Toggle crop",
	"o": "Toggle vigenette overlay",
	"s": "Go to sleep",
	"w": "Wake up",
	"0": "Go to position before tube 1",
	"1-7": "Collect from test tube n",
	" ": "Dispense {}".format(DISPENSE_uL),
	"lmb": "Select IK target",
	"u": "Uncalibrate the motors, freeing up movement if there's been a slip",
	"d": "open drain",
	"p": "close/plug drain",
}

def print_help_text():
	print("~~ HELP ~~")
	for k in helps:
		print("'{}':\t{}".format(k, helps[k]))
	print()

class Interface(Window):
	def __init__(self):
		super().__init__()

		self.target_x_rel = 0
		self.target_y_rel = 0

		self.do_crop = True
		self.do_mask = False

		# tcp is default
		# top-cam-crop not working with ffmpeg pipelines for some reason. And gstreamer isn't installed
		self.stream = cv2.VideoCapture("rtsp://DEPTH:8554/top-cam")#, cv2.CAP_FFMPEG)

		self.crop_config = read_remote_crop_config()
		if self.crop_config is None:
			print("no crop config, quitting")
			self.exit()
		self.crop_mag = self.crop_config['right_abs'] - self.crop_config['left_abs']

		self.load_mask()

	def load_mask(self):
		mask = cv2.imread(TOP_MASK)
		if len(mask) < 1:
			print("mask not loaded, quitting")
			self.exit()
			return
		
		resized_mask = cv2.resize(mask, (self.crop_mag, self.crop_mag))
		self.mask = 1 - resized_mask[:,:,0] / 255.0

	# go from absolute pixels to (-1, 1) relative coordinates
	def abs_to_rel(self, x, y):
		x -= self.crop_mag / 2.0
		y -= self.crop_mag / 2.0

		if not self.do_crop:
			x -= self.crop_config['left_abs']
			y -= self.crop_config['top_abs']

		# scale to make it (-1, +1)
		x /= self.crop_mag / 2.0 
		y /= self.crop_mag / 2.0

		# cartesian opposite to pixel coords
		y *= -1
	
		return x, y
	
	# go from (-1, 1) relative coordinates to absolute pixels
	def rel_to_abs(self, x, y):
		# cartesian opposite to pixel coords
		y *= -1

		x *= self.crop_mag / 2.0 
		y *= self.crop_mag / 2.0

		x += self.crop_mag / 2.0
		y += self.crop_mag / 2.0

		if not self.do_crop:
			x += self.crop_config['left_abs']
			y += self.crop_config['top_abs']
		
		return int(x), int(y)

		
	def mouse_handler(self, event, x, y, flags, param):
		if flags & cv2.EVENT_FLAG_LBUTTON:
			xr, yr = self.abs_to_rel(x, y)
			# reduce to unit length
			m = abs(math.hypot(xr, yr))
			if m > 1:
				xr /= m
				yr /= m

			self.target_x_rel = xr
			self.target_y_rel = yr

			mc.goto_xy(self.target_x_rel, self.target_y_rel)
			print("goto_xy", self.target_x_rel, self.target_y_rel)


	
	def keyboard_handler(self, key):
		super().keyboard_handler(key)
		if key == ord('h'):
			print_help_text()
		if key == ord('c'):
			if self.crop_config != None:
				self.do_crop = not self.do_crop
			else:
				print("cannot crop because config could not be loaded")
		if key == ord('o'):
			if len(self.mask) > 0:
				self.do_mask = not self.do_mask
			else:
				print("cannot mask because no mask loaded")
		if key == ord(' '):
			print("dispense at", self.target_x_rel, self.target_y_rel)
			mc.dispense(DISPENSE_uL)
		if key == ord('s'):
			print("sent sleep request")
			mc.sleep()
		if key == ord('w'):
			print("sent wake request")
			mc.wake()
		if key >= ord('0') and key <= ord('7'):
			num = key - ord('0')
			print("selected position", num)
			mc.collect(num)
		if key == ord('u'):
			print("uncalibrating")
			mc.uncalibrate()
		if key == ord('d'):
			print("open drain...")
			mc.set_drain(True)
		if key == ord('p'):
			print("close drain...")
			mc.set_drain(False)


	def crop(self, frame):
		top = self.crop_config['top_abs']
		bottom = self.crop_config['bottom_abs']
		left = self.crop_config['left_abs']
		right = self.crop_config['right_abs']

		return frame[top:bottom, left:right]

	def update(self):
		ret, frame = self.stream.read()
		if ret == False:
			print("Frame empty, exiting")
			self.exit()

		if self.do_crop:
			frame = self.crop(frame)

		if self.do_mask:
			start = time.time()
			if self.do_crop:
				image.overlay_image_alpha(frame, np.zeros((self.crop_mag, self.crop_mag, 3)), 0, 0, self.mask)
			else:
				image.overlay_image_alpha(frame, np.zeros((self.crop_mag, self.crop_mag, 3)), self.crop_config['left_abs'], self.crop_config['top_abs'], self.mask)
			print(time.time() - start)
		


		cv2.circle(frame,self.rel_to_abs(self.target_x_rel,self.target_y_rel),10,(0,0,255),2, cv2.LINE_AA)

		return frame

if __name__ == "__main__":
	win = Interface()
	win.loop()