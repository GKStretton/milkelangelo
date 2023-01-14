# This is an automated system for editing the video content of A Study of Light
import argparse
import os
from moviepy.editor import *
import moviepy.video.VideoClip as VideoClip
import moviepy.video.io.VideoFileClip as VideoFileClip
import yaml
import pycommon.machine_pb2 as pb
import typing
import subprocess
from datetime import datetime

TOP_CAM = "top-cam"
FRONT_CAM = "front-cam"

class CropConfig:
	def __init__(self, path):
		self.is_loaded = False
		yml = self._load_crop_config(path)
		if yml is None:
			return
		
		self.is_loaded = True
		self.raw_yml = yml
		self.x1 = yml['left_abs']
		self.x2 = yml['right_abs']
		self.y1 = yml['top_abs']
		self.y2 = yml['bottom_abs']

	def _load_crop_config(self, path):
		config = None
		try:
			with open(path, 'r') as f:
				config = yaml.load(f, Loader=yaml.FullLoader)
		except FileNotFoundError as err:
			print("error loading crop config at '{}': {}", path, err)

		return config

class FootagePiece:
	def __init__(self, path):
		# create VideoClip 
		self.video = VideoFileClip.VideoFileClip(path)
		# load crop configs 
		self.crop_config = CropConfig(path + ".yml")
		# calculate absolute start timestamp 
		result = subprocess.run(["./user-tools/get-creation-timestamp.sh", path], stdout=subprocess.PIPE)
		result.check_returncode()
		raw_ts = result.stdout.decode("utf-8")
		self.creation_timestamp = float(raw_ts)
		print("clip creation date:",datetime.utcfromtimestamp(self.creation_timestamp).strftime('%Y-%m-%d %H:%M:%S'))

	def get_clip(self) -> VideoClip.VideoClip:
		return self.video
	
	def get_crop_config(self) -> CropConfig:
		if self.crop_config.is_loaded:
			return self.crop_config
		return None
	
	def get_creation_timestamp(self) -> float:
		return self.creation_timestamp

# FootageWrapper abstracts out any separate video recordings from paused sessions
# and lets us get subclips by absolute timestamps rather than file-relative.
# it will return crop information with the subclips.
class FootageWrapper:
	def __init__(self, *paths):
		self.clips = []
		for p in paths:
			self.clips.append(FootagePiece(p))

	def get_subclip(camera: str, start_t: float, end_t: float) -> typing.Tuple[VideoClip.VideoClip, CropConfig]:
		pass
	
if __name__ == "__main__":
	parser = argparse.ArgumentParser()
	parser.add_argument("-d", "--session-dir", action="store", help="session content directory")
	parser.add_argument("-i", "--inspect", action="store_true", help="If true, detailed information will be shown on the video")
	args = parser.parse_args()

	print(args.session_dir)

	state_reports_path = os.path.join(args.session_dir, "state-reports.yml")
	with open(state_reports_path, 'r') as f:
		state_reports = yaml.load(f, yaml.FullLoader)
	
	print(len(state_reports))
	print(state_reports[0])

	filepath = os.path.join(args.session_dir, "video/raw/top-cam/1.mp4")
	clip = FootagePiece(filepath)

	print(clip.creation_timestamp)
	print(clip.get_clip().duration)
	print(clip.crop_config.raw_yml)