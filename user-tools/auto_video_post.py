# This is an automated system for editing the video content of A Study of Light
# if this becomes too slow, worth looking at using direct [ffmpeg bindings](https://github.com/kkroening/ffmpeg-python)
import argparse
import typing
import os
from moviepy.editor import *
import moviepy.video.VideoClip as VideoClip
import moviepy.video.io.VideoFileClip as VideoFileClip
import yaml
import pycommon.machine_pb2 as pb
from google.protobuf.json_format import ParseDict
from pycommon.footage import FootageWrapper
import pycommon.util as util

TOP_CAM = "top-cam"
FRONT_CAM = "front-cam"

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

# Scenes for defining composition of top and front cams
SCENE_UNDEFINED = 0
SCENE_FRONT_ONLY = 1
SCENE_FRONT_PRIMARY = 2
SCENE_TOP_PRIMARY = 3
SCENE_TOP_ONLY = 4

FORMAT_UNDEFINED = 0
FORMAT_LANDSCAPE = 1
FORMAT_PORTRAIT = 2

class ContentGenerator:
	def __init__(self, args: argparse.Namespace):
		self.args = args

		## Gather resources
		print("~~~ GATHERING RESOURCES ~~~")

		self.session_metadata = get_session_metadata(args)
		self.state_reports = get_state_reports(args)

		# load camera footage
		content_path = get_session_content_path(args)
		
		self.top_footage = FootageWrapper(os.path.join(content_path, "video/raw/" + TOP_CAM))
		self.front_footage = FootageWrapper(os.path.join(content_path, "video/raw/" + FRONT_CAM))

		print()

	def generate_content(self, args: argparse.Namespace, content_type: str):
		print("Iterating state reports...")
		subclips: typing.List[VideoClip.VideoClip] = []

		section_start_ts = 0
		section_properties = {
			"scene": SCENE_UNDEFINED,
			"speed": 1.0,
		}

		# content already generated up to this point, so skip any state reports without info past this point
		already_generated_up_to = 0

		for i in range(len(self.state_reports)):
			report = ParseDict(self.state_reports[i], pb.StateReport())

			status_name = pb.Status.Name(report.status)
			ts = float(report.timestamp_unix_micros) / 1.0e6

			print("{}\t{}     ({:.2f})\t{}".format(i, util.ts_format(ts), ts, status_name))


			#! outline:
			# get new_section_properties under this state report
			# if new_section_properties is different to section_properties, begin a new section
				# generate_section on the previous properties, from section_start_ts to latest ts.
				# set section_start_ts to latest ts
				# set section_properties to new ones
			# if dispense, immediately generate the section and mark "already_generated_up_to"

			# ! end outline

			# # clip interval
			# start_ts = ts
			# end_ts = ts + 3600*24*365 # 1 year in future, guaranteed to be at end of clip...

			# # If there is an upcoming state report, set that as end_ts instead
			# if i + 1 < len(self.state_reports):
			# 	next_report = ParseDict(self.state_reports[i+1], pb.StateReport())
			# 	end_ts = float(next_report.timestamp_unix_micros) / 1.0e6
			
			# print("\tInterval range: {:.2f} -> {:.2f}\t({:.2f})".format(start_ts, end_ts, end_ts-start_ts))


		final_clip = concatenate_videoclips(subclips)

		self.write_video(content_type, final_clip)
	
	# returns properties for this section. if the second parameter is not 0, 
	# this is a "forced_duration". a forced duration requires these properties be
	# maintained for this time, even if the state reports change.
	def get_section_properties(self, state_report) -> typing.Tuple[properties, float]:
		pass
	

	# generates a section of the content, one subclip.
	def generate_section(self, start_ts: float, end_ts: float, scene: int, speed: float = 1.0, format: int = FORMAT_LANDSCAPE) -> VideoClip:
		# TOP-CAM
		print("\tGetting top-cam clip...")
		top_clip, top_crop = self.top_footage.get_subclip(start_t=start_ts, end_t=end_ts)
		if top_clip is None:
			print("\tNo footage of top-cam, skipping")
			return
		print("\ttop-cam footage duration:\t{:.2f}".format(top_clip.duration))


		# FRONT-CAM
		print("\tGetting front-cam clip...")
		front_clip, front_crop = self.front_footage.get_subclip(start_t=start_ts, end_t=end_ts)
		if front_clip is None:
			print("\tNo footage of front-cam, skipping")
			return
		print("\tfront-cam footage duration:\t{:.2f}".format(front_clip.duration))

		# CROP
		# todo: crop

		# SPEED
		if speed != 1:
			top_clip = top_clip.speedx(speed)
			front_clip = front_clip.speedx(speed)
		
		landscape_dim = (1920, 1080)
		portrait_dim = (1080, 1920)

		if format == FORMAT_LANDSCAPE:
			if scene == SCENE_FRONT_PRIMARY:
				top_clip = top_clip.resize(0.5).set_position((50, 50))
				return CompositeVideoClip([front_clip, top_clip], size=landscape_dim)
			elif scene == SCENE_TOP_PRIMARY:
				front_clip = front_clip.resize(0.5).set_position((50, 50))
				return CompositeVideoClip([top_clip, front_clip], size=landscape_dim)
			else:
				print("scene {} not supported".format(scene))
				return None
		elif format == FORMAT_PORTRAIT:
			if scene == SCENE_FRONT_PRIMARY or scene == SCENE_TOP_PRIMARY:
				top_clip = top_clip.resize(0.5).set_position((100, 100))
				front_clip = front_clip.resize(0.5).set_position((100, 1000))
				return CompositeVideoClip([front_clip, top_clip], size=portrait_dim)
			else:
				print("scene {} not supported".format(scene))
				return None
		else:
			print("format {} not supported, returning".format(format))
			return None

	def write_final_clip(self, content_type: str, clip):
		# Get filename for writing
		output_dir = os.path.join(get_session_content_path(self.args), "video/post/")
		if not os.path.exists(output_dir):
			os.mkdir(output_dir)
		output_file = os.path.join(output_dir, content_type + ".mp4")
		i = 1
		while os.path.exists(output_file):
			output_file = os.path.join(output_dir, "{}.{}.mp4".format(content_type, i))
			i+=1

		clip.write_videofile(output_file, codec='libx264', fps=30)

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

	cg.generate_content(args, "test")
	# cg.generate_content(LONGFORM_1)
	# cg.generate_content(SHORTFORM_FULL)
	# cg.generate_content(SHORTFORM_HIGHLIGHTS)