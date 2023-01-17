import moviepy.video.VideoClip as VideoClip
import moviepy.video.io.VideoFileClip as VideoFileClip
import typing
import pycommon.util as util
import os
from pycommon.crop_util import CropConfig

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
		return util.ts_format(self.start_timestamp)
	
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