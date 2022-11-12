import cv2
from pycommon.window import Window

class Interface(Window):
	def __init__(self):
		super().__init__()

		self.target_x = 0
		self.target_y = 0

		# tcp is default
		# top-cam-crop not working with ffmpeg pipelines for some reason. And gstreamer isn't installed
		self.stream = cv2.VideoCapture("rtsp://DEPTH:8554/top-cam")#, cv2.CAP_FFMPEG)

	def mouse_handler(self, event, x, y, flags, param):
		if flags & cv2.EVENT_FLAG_LBUTTON:
			self.target_x = x
			self.target_y = y

	def update(self):
		ret, frame = self.stream.read()
		if ret == False:
			print("Frame empty, exiting")
			self.exit()

		cv2.circle(frame,(self.target_x,self.target_y),10,(0,0,255),2, cv2.LINE_AA)

		return frame

if __name__ == "__main__":
	win = Interface()
	win.loop()