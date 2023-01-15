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
		self.raw_yml = yml
		if yml is None:
			return
		
		self.is_loaded = True
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
		print("\tLoading FootagePiece:", path)

		# create VideoClip 
		self.video = VideoFileClip.VideoFileClip(path)
		# load crop configs 
		self.crop_config = CropConfig(path + ".yml")

		# calculate absolute start timestamp
		with open(path + ".creationtime", 'r') as f:
			unixtime = f.readline()

			# timestamp unix in seconds with decimal
			self.creation_timestamp = float(unixtime)
		
		print("\t\tcc:\t\t", self.crop_config.raw_yml)
		print("\t\tstart:\t\t", self.get_creation_timestamp_string())
		print("\t\tduration:\t {}s".format(self.video.duration))

	def get_clip(self) -> VideoClip.VideoClip:
		return self.video
	
	def get_crop_config(self) -> CropConfig:
		if self.crop_config.is_loaded:
			return self.crop_config
		return None
	
	def get_creation_timestamp(self) -> float:
		return self.creation_timestamp

	def get_creation_timestamp_string(self) -> str:
		return datetime.utcfromtimestamp(self.creation_timestamp).strftime('%Y-%m-%d %H:%M:%S')


# FootageWrapper abstracts out any separate video recordings from paused sessions
# and lets us get subclips by absolute timestamps rather than file-relative.
# it will return crop information with the subclips.
class FootageWrapper:
	def __init__(self, footagePath):
		print("Loading footage from directory:", footagePath)
		self.clips = []
		for file in os.listdir(footagePath):
			# get every .mp4 in the directory
			if file.endswith(".mp4"):
				# create a FootagePiece for each
				path = os.path.join(footagePath, file)
				self.clips.append(FootagePiece(path))

	def get_subclip(start_t: float, end_t: float) -> typing.Tuple[VideoClip.VideoClip, CropConfig]:
		#? how to handle transfer between clips?
		pass

def get_session_metadata(args: argparse.Namespace):
	filename = "{}_session.yml".format(args.session_number)
	path = os.path.join(args.base_dir, "session_metadata", filename)
	yml = None
	with open(path, 'r') as f:
		yml = yaml.load(f, Loader=yaml.FullLoader)
	print("Loaded session metadata:", yml)
	return yml
	
def get_session_content_path(args: argparse.Namespace):
	return os.path.join(args.base_dir, "session_content", args.session_number)

def get_state_reports(args: argparse.Namespace):
	content_path = get_session_content_path(args)
	state_reports = None
	with open(os.path.join(content_path, "state-reports.yml"), 'r') as f:
		state_reports = yaml.load(f, yaml.FullLoader)
	print("Loaded {} state report entries".format(len(state_reports)))
	return state_reports

def get_cam_footage_wrapper(args: argparse.Namespace, cam: str):
	content_path = get_session_content_path(args)
	return FootageWrapper(os.path.join(content_path, "video/raw/" + cam))


class ContentGenerator:
	def __init__(self, args: argparse.Namespace):
		## Gather resources
		print("~~~ GATHERING RESOURCES ~~~")

		self.session_metadata = get_session_metadata(args)
		self.state_reports = get_state_reports(args)

		# load camera footage
		self.top_footage = get_cam_footage_wrapper(args, TOP_CAM)
		self.front_footage = get_cam_footage_wrapper(args, FRONT_CAM)

		print()

	def generate_content(self, content_type: str):
		pass



if __name__ == "__main__":
	parser = argparse.ArgumentParser()
	parser.add_argument("-d", "--base-dir", action="store", help="base directory containing session_content and session_metadata")
	parser.add_argument("-n", "--session-number", action="store", help="session number e.g. 5")
	parser.add_argument("-i", "--inspect", action="store_true", help="If true, detailed information will be shown on the video")
	args = parser.parse_args()
	print("Launching auto_video_post for session {} in '{}'\n".format(args.session_number, args.base_dir))

	cg = ContentGenerator(args)

	# cg.generate_content(LONGFORM_1)
	# cg.generate_content(SHORTFORM_FULL)
	# cg.generate_content(SHORTFORM_HIGHLIGHTS)
