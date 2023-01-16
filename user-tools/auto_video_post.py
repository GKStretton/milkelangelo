# This is an automated system for editing the video content of A Study of Light
import argparse
import os
from moviepy.editor import *
import moviepy.video.VideoClip as VideoClip
import moviepy.video.io.VideoFileClip as VideoFileClip
import yaml
import pycommon.machine_pb2 as pb
from google.protobuf.json_format import Parse, ParseDict
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


# FootagePiece handles loading and storage of a video clip, its crop config, and
# timestamps. It supports getting subclips by absolute timestamp, wrapping
# moviepy's file-relative timestamps.
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
			self.start_timestamp = float(unixtime)
		
		print("\t\tcc:\t\t", self.crop_config.raw_yml)
		print("\t\tstart:\t\t {} ({})".format(self.get_start_timestamp(), self.get_start_timestamp_string()))
		print("\t\tduration:\t {}".format(self.video.duration))

	def get_clip(self) -> VideoClip.VideoClip:
		return self.video
	
	def get_crop_config(self) -> CropConfig:
		if self.crop_config.is_loaded:
			return self.crop_config
		return None
	
	def get_start_timestamp(self) -> float:
		return self.start_timestamp

	def get_start_timestamp_string(self) -> str:
		return datetime.utcfromtimestamp(self.start_timestamp).strftime('%Y-%m-%d %H:%M:%S')
	
	def get_end_timestamp(self) -> float:
		return self.start_timestamp + self.video.duration
	
	def contains_timestamp(self, t: float):
		return t >= self.get_start_timestamp() and t <= self.get_end_timestamp()
	
	# this returns a subclip within decimal absolute unix timestamps. Handles the
	# case where end_t is outside the footage by setting end to footage end
	def get_subclip_from_timestamps(self, start_t: float, end_t: float) -> VideoClip.VideoClip:
		if not self.contains_timestamp(start_t):
			return None
		
		start_relative = start_t - self.start_timestamp
		end_relative = end_t - self.get_end_timestamp()
		if end_relative > self.video.duration:
			end_relative = self.video.duration

		return self.video.subclip(start_relative, end_relative)


# FootageWrapper abstracts out any separate video recordings from paused sessions
# it will return crop information with the subclips.
class FootageWrapper:
	def __init__(self, footagePath):
		print("Loading footage from directory:", footagePath)
		self.clips: typing.List[FootagePiece] = []
		for file in os.listdir(footagePath):
			# get every .mp4 in the directory
			if file.endswith(".mp4"):
				# create a FootagePiece for each
				path = os.path.join(footagePath, file)
				self.clips.append(FootagePiece(path))

	def get_subclip(self, start_t: float, end_t: float) -> typing.Tuple[VideoClip.VideoClip, CropConfig]:
		# if start_t is in clip x, we ignore everything after clip x. So each
		# state report interval can only be in one clip. This is okay because
		# on start / resume, a state report will be triggered.
		for _, v in enumerate(self.clips):
			if v.contains_timestamp(start_t):
				# this clip contains footage of the state report
				return v.get_subclip_from_timestamps(start_t, end_t), v.get_crop_config()
		
		# no clip with footage
		return None, None
	
	def test(self):
		start = self.clips[0].start_timestamp + 8
		end = start + 3.5
		c, _ = self.get_subclip(start, end)
		c.preview()


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
		print("Iterating state reports...")
		for i, sr in enumerate(self.state_reports):
			report = ParseDict(sr, pb.StateReport())
			print("{}\t - {}".format(i, report.status))

			# todo: another test session, with all testable features
			# todo: verify it here
			# todo: generate a video clip that is normal, except with overlay of state report information

	def test(self):
		self.top_footage.test()



if __name__ == "__main__":
	parser = argparse.ArgumentParser()
	parser.add_argument("-d", "--base-dir", action="store", help="base directory containing session_content and session_metadata")
	parser.add_argument("-n", "--session-number", action="store", help="session number e.g. 5")
	parser.add_argument("-i", "--inspect", action="store_true", help="If true, detailed information will be shown on the video")
	args = parser.parse_args()
	print("Launching auto_video_post for session {} in '{}'\n".format(args.session_number, args.base_dir))

	cg = ContentGenerator(args)

	# cg.test()

	cg.generate_content("")
	# cg.generate_content(LONGFORM_1)
	# cg.generate_content(SHORTFORM_FULL)
	# cg.generate_content(SHORTFORM_HIGHLIGHTS)