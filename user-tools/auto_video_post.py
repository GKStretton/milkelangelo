# This is an automated system for editing the video content of A Study of Light
# if this becomes too slow, try using nvenc as moviepy ffmpeg codec
# alternately, at using direct [ffmpeg bindings](https://github.com/kkroening/ffmpeg-python)
import argparse
import os

from videoediting.content_generator import ContentGenerator
from videoediting.constants import *
from videoediting.footage import FootageWrapper
import videoediting.loaders as loaders
from google.protobuf.json_format import ParseDict
import pycommon.machine_pb2 as pb
import videoediting.section_properties as properties

from dataclasses import dataclass
from enum import Enum
from abc import ABC, abstractmethod

@dataclass
class SectionProperties:
	speed: float = 1.0

# This is a descriptor (list of timestamps and video properties) for a single
# piece of content (video)
class ContentDescriptor():
	def set_state_report(self, timestamp: float, state_report):
		pass

	def set_properties(self, timestamp: float, properties: SectionProperties):
		# todo: only add if it's changed
		pass

if __name__ == "__main__":
	parser = argparse.ArgumentParser()
	parser.add_argument("-d", "--base-dir", action="store", help="base directory containing session_content and session_metadata")
	parser.add_argument("-n", "--session-number", action="store", help="session number e.g. 5")
	parser.add_argument("-i", "--inspect", action="store_true", help="If true, detailed information will be shown on the video")
	parser.add_argument("-x", "--dry-run", action="store_true", help="If true, this will be a dry run and no content will be generated")
	parser.add_argument("-e", "--end-at", action="store", help="If set, content will be ended after this timestamp (s)")
	parser.add_argument("-p", "--preview", action="store_true", help="If true, final video will be previewed rather than written")
	parser.add_argument("-s", "--show", action="store", help="If true, show frame at this timestamp (s)")

	args = parser.parse_args()
	print(f"Launching auto_video_post for session {args.session_number} in '{args.base_dir}'\n")

	session_metadata = loaders.get_session_metadata(args)
	state_reports = loaders.get_state_reports(args)

	# load camera footage
	content_path = loaders.get_session_content_path(args)
	top_footage = FootageWrapper(os.path.join(content_path, "video/raw/" + TOP_CAM))
	front_footage = FootageWrapper(os.path.join(content_path, "video/raw/" + FRONT_CAM))

	propertyList = {}
	state = {}
	descriptor = ContentDescriptor()
	for i in range(len(state_reports)):
		report = ParseDict(state_reports[i], pb.StateReport())
		report_ts = float(report.timestamp_unix_micros) / 1.0e6

		props = properties.get_section_properties(state, report, TYPE_LONGFORM)
		descriptor.set_state_report(report_ts, report)
		descriptor.set_properties(report_ts, props)
	
	# generate overlay from state reports
	# generate overlay from properties
	# composite the above 2

	# > overlay track
	# > content track

	# todo: make footage wrapper support pre- and-post padding rather than clipping.

	# iterate properties in content descriptor
	# append subclips for overlay and content accordingly

	# launch preview application, or save