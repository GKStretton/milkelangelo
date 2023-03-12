import moviepy.video.VideoClip as VideoClip
from moviepy.editor import concatenate_videoclips
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
		if start_t < self.get_end_timestamp():
			if end_t is None or end_t > self.get_start_timestamp():
				return True
		return False
	
	# this returns a subclip within decimal absolute unix timestamps. 
	# Return any footage from this video in the stated range.
	# (footage, absolute_start, absolute_end)
	def get_subclip_from_timestamps(self, start_t: float, end_t: float) -> typing.Tuple[VideoClip.VideoClip, float, float]:
		if not self.intersects_timestamp_range(start_t, end_t):
			return None, None, None

		if end_t is None:
			end_t = self.get_end_timestamp()

		
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
		return subclip, absolute_start_t, absolute_end_t


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
		print()
	
	def get_start_timestamp(self):
		if len(self.clips) == 0:
			return None
		return self.clips[0].get_start_timestamp()

	# get subclip by absolute timestamp range, including padding
	def get_subclip(self, start_t: float, end_t: float) -> typing.Tuple[VideoClip.VideoClip, CropConfig]:
		if len(self.clips) == 0:
			print("get_subclip found self.clips to be empty")
			exit(1)

		# clips with absolute timestamp ranges
		# [(clip, start_a, end_a)]
		subclips = []
		crop_config = None
		for _, v in enumerate(self.clips):
			clip, start_a, end_a = v.get_subclip_from_timestamps(start_t, end_t)
			if clip is not None:
				# this clip contains footage of the state report
				subclips.append((clip, start_a, end_a))
				crop_config = v.get_crop_config()
				print(f"\tFound footage from {start_a} to {end_a} ({end_a-start_a})")
		
		# padding template
		s = self.clips[0].get_clip().size
		dur = end_t-start_t if end_t is not None else 2
		padding = VideoClip.ColorClip(s, color=(0, 0, 0), duration=dur)
		
		# generate black video with correct duration if there's no footage
		if len(subclips) == 0:
			return padding, crop_config
		
		# if there is footage...
		subclips_with_padding = []

		# add any pre-video padding
		initial_pad_duration = subclips[0][1] - start_t
		if initial_pad_duration > 0:
			print(f"\tAdding {initial_pad_duration:.2f} of initial padding")
			subclips_with_padding.append(padding.set_duration(initial_pad_duration))
		
		# iterate subclips
		for i in range(len(subclips)):
			# add this clip
			subclips_with_padding.append(subclips[i][0])
			# if there's another one next, 
			if i < len(subclips) - 1:
				# add padding until the start of next one
				padding_duration = subclips[i+1][1] - subclips[i][2]
				subclips_with_padding.append(padding.set_duration(padding_duration))
		
		# add end padding
		if end_t is not None:
			end_padding_duration = max(0, end_t - subclips[-1][2])
			if end_padding_duration > 0:
				subclips_with_padding.append(padding.set_duration(end_padding_duration))

		return concatenate_videoclips(subclips_with_padding), crop_config
	
	def test(self):
		start = self.clips[0].start_timestamp + 8
		end = start + 3.5
		c, _ = self.get_subclip(start, end)
		c.preview()