import moviepy.video.VideoClip as VideoClip
import moviepy.video.io.VideoFileClip as VideoFileClip
import time
import typing
import pycommon.util as util
import os
from pycommon.crop_util import CropConfig

# FootagePiece handles loading and storage of a video clip, its crop config, and
# timestamps. It supports getting subclips by absolute timestamp, wrapping
# moviepy's file-relative timestamps.
class FootagePiece:
	def __init__(self, path):
		self.file_name = os.path.basename(path)
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
		
		print("\t\tcc:\t\t x1={}; x2={}; y1={}; y2={}".format(self.crop_config.x1, self.crop_config.x2, self.crop_config.y1, self.crop_config.y2))
		print("\t\trange:\t\t {:.2f} - {:.2f}".format(self.get_start_timestamp(), self.get_end_timestamp()))
		print("\t\tduration:\t {:.2f}".format(self.video.duration))

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
	
	# returns true if the video absolute time range overlaps with the given
	# absolute time range. i.e. does this video contain any content from this
	# time range?
	def intersects_timestamp_range(self, start_t: float, end_t: float):
		if start_t >= self.get_start_timestamp() and start_t <= self.get_end_timestamp():
			return True
		if end_t >= self.get_start_timestamp() and end_t <= self.get_end_timestamp():
			return True

		return False
	
	# this returns a subclip within decimal absolute unix timestamps. 
	# Return any footage from this video in the stated range.
	def get_subclip_from_timestamps(self, start_t: float, end_t: float) -> VideoClip.VideoClip:
		if not self.intersects_timestamp_range(start_t, end_t):
			return None
		
		print("\t\tGathering subclip for range ({:.2f}, {:.2f}) in {}".format(start_t, end_t, self.file_name))
		
		# video confined to range inside both video and stated range
		absolute_start_t = max(start_t, self.get_start_timestamp())
		absolute_end_t = min(end_t, self.get_end_timestamp())

		print("\t\tApplicable absolute range considered ({:.2f}, {:.2f}) in {}".format(absolute_start_t, absolute_end_t, self.file_name))

		start_relative = absolute_start_t - self.get_start_timestamp()
		end_relative = absolute_end_t - self.get_start_timestamp()
		if start_relative < 0:
			print("\t\tstart_relative ({}) < 0 in {}, exiting".format(start_relative, self.file_name))
			exit(1)
		if end_relative > self.video.duration:
			end_relative = self.video.duration
			# should only be very small, floating point inaccuracies
			if end_relative - self.video.duration > 0.01:
				print("\t\tend_relative {} significantly bigger than video duration {} in {}, exiting".format(end_relative, self.video.duration, self.file_name))

		print("\t\tApplicable relative range considered ({:.2f}, {:.2f}) in {}".format(start_relative, end_relative, self.file_name))

		before = time.time()
		subclip = self.video.subclip(start_relative, end_relative)
		print("\t\tvideo.subclip took {}s".format(time.time() - before))
		return subclip


# FootageWrapper abstracts out any separate video recordings from paused sessions
# it will return crop information with the subclips.
class FootageWrapper:
	def __init__(self, footagePath):
		print("Loading footage from directory:", footagePath)
		self.clips: typing.List[FootagePiece] = []
		for file in sorted(os.listdir(footagePath)):
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
			clip = v.get_subclip_from_timestamps(start_t, end_t)
			if clip is not None:
				# this clip contains footage of the state report
				return clip, v.get_crop_config()
			else:
				continue
		
		# no clip with footage
		return None, None
	
	def test(self):
		start = self.clips[0].start_timestamp + 8
		end = start + 3.5
		c, _ = self.get_subclip(start, end)
		c.preview()