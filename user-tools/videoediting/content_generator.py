import argparse
import os
import typing
from termcolor import colored
import cv2 as cv

from moviepy.editor import *
import moviepy.video.VideoClip as VideoClip
import moviepy.video.fx.all as vfx

from google.protobuf.json_format import ParseDict
from google.protobuf.json_format import MessageToJson
import pycommon.machine_pb2 as pb
import pycommon.util as util
import pycommon.image as image

from videoediting.constants import *
from videoediting.footage import FootageWrapper
import videoediting.loaders as loaders
import videoediting.section_properties as properties

from dataclasses import dataclass
from enum import Enum
from abc import ABC, abstractmethod

class FootageType(Enum):
	TOP_CAM = 1
	FRONT_CAM = 2



# Builds the properties that will be used to render a specific video format
# this is created by the user of this framework
class ContentDescriptionsBuilder(ABC):
	# the function that turns state reports or video into ContentDescriptor(s)
	@abstractmethod
	def build_content_descriptors(
		self,
		state_reports,
		footage_optional: typing.Dict[FootageType, FootageWrapper]
	) -> typing.List[ContentDescriptor]:
		pass

class ContentGenerator:
	def __init__(self, args: argparse.Namespace):
		self.args = args

		## Gather resources
		print("~~~ GATHERING RESOURCES ~~~")

		self.session_metadata = loaders.get_session_metadata(args.base_dir, args.session_number)
		self.state_reports = loaders.get_state_reports(args)

		if not args.dry_run:
			# load camera footage
			content_path = loaders.get_session_content_path(args)
			
			self.top_footage = FootageWrapper(os.path.join(content_path, "video/raw/" + TOP_CAM))
			self.front_footage = FootageWrapper(os.path.join(content_path, "video/raw/" + FRONT_CAM))

		print()
	
	def generate_content(self, content_builder: ContentDescriptionsBuilder):
		# ! ContentDescriptor
		# props is list of timestamped properties:
		# [timestamp: {
			# state_report,
			# video_properties,
		# }
		# , ...]
		# and has some global properties:
		# [intro_still_settings, outro_settings, session_meta]
		descriptors = content_builder.build_content_descriptors()

		# list of ContentPiece objects
		# content_list = [(footage, info_overlay, information, eventList/stateReports?), ...]
		# content_list = generate_footage_from_properties_list(props...)

		# return footage, info_overlay
	