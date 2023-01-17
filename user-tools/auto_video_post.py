# This is an automated system for editing the video content of A Study of Light
import argparse
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


class ContentGenerator:
	def __init__(self, args: argparse.Namespace):
		## Gather resources
		print("~~~ GATHERING RESOURCES ~~~")

		self.session_metadata = get_session_metadata(args)
		self.state_reports = get_state_reports(args)

		# load camera footage
		content_path = get_session_content_path(args)
		
		self.top_footage = FootageWrapper(os.path.join(content_path, "video/raw/" + TOP_CAM))
		self.front_footage = FootageWrapper(os.path.join(content_path, "video/raw/" + FRONT_CAM))

		print()

	def generate_content(self, content_type: str):
		print("Iterating state reports...")
		for i in range(len(self.state_reports)):
			report = ParseDict(self.state_reports[i], pb.StateReport())

			status_name = pb.Status.Name(report.status)
			ts = float(report.timestamp_unix_micros) / 1.0e6

			print("{}\t{}     ({})\t{}".format(i, util.ts_format(ts), ts, status_name))

			# clip interval
			start_ts = ts
			end_ts = ts + 3600*24*365 # 1 year in future, guaranteed to be at end of clip...

			# If there is an upcoming state report, set that as end_ts instead
			if i + 1 < len(self.state_reports):
				next_report = ParseDict(self.state_reports[i+1], pb.StateReport())
				end_ts = float(next_report.timestamp_unix_micros) / 1.0e6
			
			print("\tInterval range: {} -> {}\t({})".format(start_ts, end_ts, end_ts-start_ts))

			#todo: handle first state report not being in 1.mp4. Currently it will choose 2.mp4 rather
			#todo: than best of 1.

			#todo: check if pausing is reflected in the below?

			#! state reports only for 1.mp4???

			top_clip, top_crop = self.top_footage.get_subclip(start_t=start_ts, end_t=end_ts)
			print("\tTop Duration:\t{}".format(top_clip.duration))

			front_clip, front_crop = self.top_footage.get_subclip(start_t=start_ts, end_t=end_ts)
			print("\tFront Duration:\t{}".format(front_clip.duration))

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