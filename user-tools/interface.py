import cv2
import time
import math
from pycommon.window import Window
from pycommon.config_manager_client import read_remote_crop_config
import pycommon.image as image
import numpy as np
import pycommon.mqtt_client as mc

TOP_MASK = "resources/static_img/top-mask.png"
# how many microlitres to dispense at a time
DISPENSE_uL = 10
collection_volume = 30.0
# how many ml of milk / water to dispense
FLUID_ML = 200.0

STREAM=True
DO_CROP=False
DO_MASK=False
SHOW_GRID=True

helps = {
	" ": "Dispense {}".format(DISPENSE_uL),
	"b": "begin session",
	"c": "Toggle crop",
	"d": "open drain",
	"e": "end session",
	"h": "Print this help text",
	"k": "Kill / Sleep",
	"m": "toggle manual mode",
	"n": "Go to navigation node",
	"o": "Toggle vigenette overlay",
	"p": "close/plug drain",
	"s": "Shutdown cleanly",
	"v": "DRAIN fluid request",
	"u": "Uncalibrate the motors, freeing up movement if there's been a slip",
	"x": "WATER fluid request",
	"y": "MILK fluid request",
	"w": "Wake up",
	"z": "pause session",
	"0": "Go to position before tube 1",
	"1-7": "Collect from test tube n",
	"lmb": "Select IK target",
	",": "resume session",
	";": "toggle grid",
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

		self.do_crop = DO_CROP
		self.do_mask = DO_MASK
		self.show_grid = SHOW_GRID


		self.crop_config = read_remote_crop_config("crop_top-cam")
		if self.crop_config is None:
			print("no crop config, disabling crop and mask")
			self.do_crop = False
			self.do_mask = False
		else:
			self.crop_mag = self.crop_config['right_abs'] - self.crop_config['left_abs']
			self.load_mask()

		if STREAM:
			# tcp is default
			# top-cam-crop not working with ffmpeg pipelines for some reason. And gstreamer isn't installed
			self.stream = cv2.VideoCapture("rtsp://DEPTH:8554/top-cam")#, cv2.CAP_FFMPEG)
			self.stream.set(cv2.CAP_PROP_BUFFERSIZE, 2)
			if not self.stream.isOpened():
				print("Error loading webcam stream, aborting.")
				self.exit()
			print("Stream opened")

	def load_mask(self):
		mask = cv2.imread(TOP_MASK, cv2.IMREAD_UNCHANGED)
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
		if event == cv2.EVENT_LBUTTONDOWN:
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
			print("sent shutdown request")
			mc.shutdown()
		if key == ord('k'):
			print("sent kill / sleep request")
			mc.sleep()
		if key == ord('w'):
			print("sent wake request")
			mc.wake()
		if key >= ord('0') and key <= ord('7'):
			num = key - ord('0')
			print("selected position", num)
			mc.collect(num, collection_volume)
		if key == ord('u'):
			print("uncalibrating")
			mc.uncalibrate()
		if key == ord('d'):
			print("open drain...")
			mc.set_drain(True)
		if key == ord('p'):
			print("close drain...")
			mc.set_drain(False)
		if key == ord('n'):
			node = input("goto-node. Specify node number (see firmware for enum): ")
			mc.goto_node(node)
		if key == ord('m'):
			print("toggling manual mode...")
			mc.toggle_manual()
		if key == ord('b'):
			print("Starting session")
			mc.pub(mc.BEGIN_SESSION, "")
		if key == ord('e'):
			print("Ending session")
			mc.pub(mc.END_SESSION, "")
		if key == ord('z'):
			print("Pausing session")
			mc.pub(mc.PAUSE_SESSION, "")
		if key == ord(','):
			print("Resuming session")
			mc.pub(mc.RESUME_SESSION, "")
		if key == ord('v'):
			print("sending fluid req: drain")
			mc.fluid_req(mc.FLUID_DRAIN, FLUID_ML)
		if key == ord('x'):
			print("sending fluid req: water")
			mc.fluid_req(mc.FLUID_WATER, FLUID_ML)
		if key == ord('y'):
			print("sending fluid req: milk")
			mc.fluid_req(mc.FLUID_MILK, FLUID_ML)
		if key == ord(';'):
			print("toggling grid")
			self.show_grid = not self.show_grid


	def crop(self, frame):
		top = self.crop_config['top_abs']
		bottom = self.crop_config['bottom_abs']
		left = self.crop_config['left_abs']
		right = self.crop_config['right_abs']

		return frame[top:bottom, left:right]
	
	def draw_grid(self, frame):
		for y in np.linspace(-1, 1, 21):
			thickness = 1
			if y == 0:
				thickness = 2
			cv2.line(frame, self.rel_to_abs(-1, y), self.rel_to_abs(1, y), (0, 255, 0), thickness, cv2.LINE_AA)
		
		for x in np.linspace(-1, 1, 21):
			thickness = 1
			if x == 0:
				thickness = 2
			cv2.line(frame, self.rel_to_abs(x, -1), self.rel_to_abs(x, 1), (0, 255, 0), thickness)

	def update(self):
		if STREAM:
			ret, frame = self.stream.read()
			if ret == False:
				print("Frame empty, exiting")
				self.exit()
		else:
			frame = np.ones((1080, 1920, 3)) * 255

		if self.do_crop:
			frame = self.crop(frame)

		if self.do_mask:
			if self.do_crop:
				image.overlay_image_alpha(frame, np.zeros((self.crop_mag, self.crop_mag, 3)), 0, 0, self.mask)
			else:
				image.overlay_image_alpha(frame, np.zeros((self.crop_mag, self.crop_mag, 3)), self.crop_config['left_abs'], self.crop_config['top_abs'], self.mask)
			
		if self.show_grid:
			self.draw_grid(frame)

		if self.crop_config is not None:
			cv2.circle(frame,self.rel_to_abs(self.target_x_rel,self.target_y_rel),3,(0,0,255),2, cv2.LINE_AA)
		

		return frame

if __name__ == "__main__":
	mc.connect()
	win = Interface()
	win.loop()