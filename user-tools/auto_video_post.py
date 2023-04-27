# This is an automated system for editing the video content of A Study of Light
# if this becomes too slow, try using nvenc as moviepy ffmpeg codec
# alternately, at using direct [ffmpeg bindings](https://github.com/kkroening/ffmpeg-python)
import argparse
from datetime import datetime
import os
import json
import dataclasses

from betterproto import Casing
from moviepy.editor import *
import moviepy.video.VideoClip as VideoClip
from moviepy.editor import concatenate_videoclips
from videoediting.constants import *
from videoediting.footage import FootageWrapper
import videoediting.loaders as loaders
import machinepb.machine as pb
from videoediting.properties.content_property_manager import *
from videoediting.properties.factory import *
import pycommon.util as util
from videoediting.compositor import compositeContentFromFootageSubclips
from videoediting.stills import add_stills
from videoediting.dispense_metadata import DispenseMetadataWrapper

from enum import Enum
from abc import ABC, abstractmethod
import typing

FPS = 60
# how much time to offset the start timestamp of the footage.
# this is basically the latency of the streaming + recording pipeline
#* if the lights come on before the state report says so, decrease this magnitude (increase number).
TOP_CAM_TIME_OFFSET=-1.18
FRONT_CAM_TIME_OFFSET=-1.18

# This is a descriptor (list of timestamps and video properties) for a single
# piece of content (video)
class ContentDescriptor:
	def __init__(self, session_metadata, content_type: str, content_fmt: Format):
		self.session_metadata = session_metadata
		self.content_type = content_type
		self.fmt = content_fmt

		# [(timestamp_s, sr), ...]
		self.state_reports: typing.List[typing.Tuple[float, pb.StateReport, VideoState]] = []
		# [(timestamp_s, props), ...]
		self.properties: typing.List[typing.Tuple[float, SectionProperties]] = []

	def set_state_report(self, timestamp: float, state_report: pb.StateReport, video_state: VideoState):
		if len(self.state_reports) > 0 and self.state_reports[-1][0] > timestamp:
			print(f"set_state_report with timestamp {timestamp}, but previously seen state report has timestamp {self.state_reports[-1][0]}")
			exit(1)
		
		self.state_reports.append((timestamp, state_report, dataclasses.replace(video_state)))

	def set_properties(self, timestamp: float, properties: SectionProperties):
		if len(self.properties) > 0 and self.properties[-1][0] > timestamp:
			print(f"set_properties with timestamp {timestamp}, but previously seen state report has timestamp {self.properties[-1][0]}")
			exit(1)
		
		# only add prop if it's different to last
		FILTER_UNCHANGED = True
		if FILTER_UNCHANGED and len(self.properties) > 0 and self.properties[-1][1] == properties:
			return

		self.properties.append((timestamp, properties))
	
	# without any skips or speed changes, no props pass
	def _generate_raw_overlay_clip(self) -> VideoClip.VideoClip:
		if len(self.properties) == 0 or len(self.state_reports) == 0:
			print("_generate_raw_overlay_clip found no properties / state_reports")
			exit(1)
		
		if self.properties[0][0] != self.state_reports[0][0]:
			print("assumed first timestamps equal, but not")
			exit(1)

		FINAL_DURATION = 2
		
		start_timestamp = self.state_reports[0][0]

		print(f"gen raw state report overlay for {len(self.state_reports)} state reports")
		# generate raw state report overlay in normal time
		sr_clips = []
		for i in range(len(self.state_reports)):
			print(i)
			timestamp, sr, video_state = self.state_reports[i]

			# show last one for this long
			duration = FINAL_DURATION
			if i < len(self.state_reports) - 1:
				next_timestamp, _, _ = self.state_reports[i+1]
				duration = next_timestamp - timestamp
			
			sr_fmt = json.dumps(
				sr.to_dict(include_default_values=True, casing=Casing.SNAKE),
				indent=4,
				sort_keys=True,
		    )
			text_str = "STATE REPORT:\n"+util.ts_format(timestamp) + "\n" + sr_fmt
			text_str += "\n\nVIDEO STATE:\n" + str(video_state) + "\n"
			txt: TextClip = TextClip(text_str, font='DejaVu-Sans-Mono', fontsize=10, color='white', align='West')
			txt = txt.set_duration(duration)

			sr_clips.append(txt)
		
		sr_full_clip = concatenate_videoclips(sr_clips)
		return sr_full_clip

	# returns (Overlay, Content)
	def generate_content_clip(self, top_footage: FootageWrapper, front_footage: FootageWrapper) -> typing.Tuple[VideoClip.VideoClip, VideoClip.VideoClip]:
		if len(self.properties) == 0:
			print("generate_content_clip found no properties / state_reports")
			exit(1)

		raw_overlay = self._generate_raw_overlay_clip()
		start_timestamp = self.properties[0][0]

		print(f"gen content for {len(self.properties)} props")
		overlay_clips = []
		content_clips = []
		for i in range(len(self.properties)):
			print(i)
			props = self.properties[i][1]
			if props.skip:
				print("skipping")
				continue

			ts1_abs = self.properties[i][0]
			ts1_rel = ts1_abs - start_timestamp
			# default to end
			ts2_abs, ts2_rel = None, None
			# if there's another prop after this
			if i < len(self.properties) - 1:
				ts2_abs = self.properties[i+1][0]
				ts2_rel = ts2_abs - start_timestamp
			
			# build overlay subclip
			overlay_raw_subclip = raw_overlay.subclip(ts1_rel, ts2_rel)
			text_str = "PROPS:\n"+util.ts_format(self.properties[i][0]) + "\n" + props.__str__()
			# commenting out this line in /etc/ImageMagick-6/policy.xml was required:
			# <policy domain="path" rights="none" pattern="@*" />
			# https://github.com/Zulko/moviepy/issues/401#issuecomment-278679961
			txt: TextClip = TextClip(text_str, font='DejaVu-Sans-Mono', fontsize=10, color='white', align='West')
			txt = txt.set_duration(overlay_raw_subclip.duration)
			overlay_subclip = clips_array([[overlay_raw_subclip], [txt]])

			# build footage subclips
			print(f"Getting top_footage between {ts1_abs} and {ts2_abs} ({util.ts_fmt(ts1_rel)} to {util.ts_fmt(ts2_rel)})")
			top_subclip, top_crop = top_footage.get_subclip(ts1_abs, ts2_abs)
			print(f"Getting front_footage between {ts1_abs} and {ts2_abs} ({util.ts_fmt(ts1_rel)} to {util.ts_fmt(ts2_rel)})")
			front_subclip, front_crop = front_footage.get_subclip(ts1_abs, ts2_abs)

			# apply speed to all
			if props.speed != 1.0:
				overlay_subclip = overlay_subclip.speedx(props.speed)
				top_subclip = top_subclip.speedx(props.speed)
				front_subclip = front_subclip.speedx(props.speed)

			# clips should all be same length unless it's the last property.
			if i != len(self.properties) - 1 and not util.floats_are_equal(0.00001, [overlay_subclip.duration, top_subclip.duration, front_subclip.duration]):
				# these should be same length if the footage has been padded correctly''
				print("processed subclips are not same duration: {} {} {}, exiting".format(overlay_subclip.duration, top_subclip.duration, front_subclip.duration))
				exit(1)
			
			content_clips.append(compositeContentFromFootageSubclips(top_subclip, top_crop, front_subclip, front_crop, props, self.fmt, self.session_metadata))
			overlay_clips.append(overlay_subclip)


		return concatenate_videoclips(overlay_clips), concatenate_videoclips(content_clips)


def save(args, overlay: VideoClip.VideoClip, content: VideoClip.VideoClip, content_type):
	output_dir = os.path.join(loaders.get_session_content_path(args), "video/post/")
	if not os.path.exists(output_dir):
		os.mkdir(output_dir)

	i = 0
	output_file = os.path.join(output_dir, "{}.{}.mp4".format(content_type.name, i))
	while os.path.exists(output_file):
		i+=1
		output_file = os.path.join(output_dir, "{}.{}.mp4".format(content_type.name, i))
	
	overlay_file = os.path.join(output_dir, "{}-overlay.{}.mp4".format(content_type.name, i))
	content_file = output_file

	overlay_render_start = datetime.now()
	overlay.write_videofile(overlay_file, codec='libx264', fps=FPS)
	print(f"overlay generation time: {str(datetime.now() - overlay_render_start)}")
	content_render_start = datetime.now()
	content.write_videofile(content_file, codec='libx264', fps=FPS)
	print(f"content generation time: {str(datetime.now() - content_render_start)}")

def test(metadata, fmt: Format, timestamp: float, top_footage: FootageWrapper, front_footage: FootageWrapper):
	top_clip, top_crop = top_footage.get_subclip(timestamp, timestamp + 1)
	front_clip, front_crop = front_footage.get_subclip(timestamp, timestamp + 1)

	props = SectionProperties(
		scene = Scene.DUAL,
		speed = 1.0,
		skip = False,
		crop = True,
		vig_overlay = True,
	)

	res = compositeContentFromFootageSubclips(top_clip, top_crop, front_clip, front_crop, props, fmt, session_metadata)

	res.resize(0.5).show(interactive=True)

if __name__ == "__main__":
	parser = argparse.ArgumentParser()
	parser.add_argument("-d", "--base-dir", action="store", help="base directory containing session_content and session_metadata", required=True)
	parser.add_argument("-n", "--session-number", action="store", help="session number e.g. 5", required=True)
	parser.add_argument("-e", "--end-at", action="store", help="If set, content will be ended after this timestamp (s)")
	parser.add_argument("-p", "--preview", action="store_true", help="If true, final video will be previewed rather than written")
	parser.add_argument("-x", "--test", action="store_true", help="If true, run test code instead of main functionality")
	parser.add_argument("-t", "--type", action="store", help="content type of output e.g. SHORTFORM | LONGFORM", default="LONGFORM")
	parser.add_argument("-s", "--start-at", action="store", help="set final clip start to this time (s), useful with preview", default=0)
	parser.add_argument("-y", "--yes", action="store_true", help="auto-yes to render confirmation", default=False)
	# parser.add_argument("-f", "--format", action="store", help="content format of output e.g. LANDSCAPE | PORTRAIT", required=True)

	args = parser.parse_args()
	print(f"Launching auto_video_post for session {args.session_number} in '{args.base_dir}'\n")

	gen_start = datetime.now()

	content_type = ContentType.__members__[args.type]
	property_manager = create_property_manager(content_type)
	content_fmt = property_manager.get_format()

	session_metadata = loaders.get_session_metadata(args.base_dir, args.session_number)
	state_reports = loaders.get_state_reports(args)
	dispense_metadata_wrapper = DispenseMetadataWrapper(args.base_dir, args.session_number)

	# load camera footage
	content_path = loaders.get_session_content_path(args)
	top_footage = FootageWrapper(os.path.join(content_path, "video/raw/" + TOP_CAM), timeOffset=TOP_CAM_TIME_OFFSET)
	front_footage = FootageWrapper(os.path.join(content_path, "video/raw/" + FRONT_CAM), timeOffset=FRONT_CAM_TIME_OFFSET)

	if args.test:
		exit(0)

	# unused
	state = VideoState()

	start_ts = state_reports[0][0]
	descriptor = ContentDescriptor(session_metadata, content_type, content_fmt)
	# no properties may be set until after this time
	generated_until = 0

	# Iterate state reports
	for i, (report_ts, report) in enumerate(state_reports):
		# Get section properties
		props, delay, min_duration = property_manager.get_section_properties(state, report, dispense_metadata_wrapper)

		# Always set state report in descriptor for use in the overlays
		descriptor.set_state_report(report_ts, report, state)

		# if we've generated beyond everything this report covers, skip it
		if (
			i+1 < len(state_reports) and
			generated_until >= state_reports[i+1][0] and
			generated_until >= report_ts + delay + min_duration
		):
			continue

		# only change properties after generated_until
		ts = max(generated_until, report_ts + delay)
		descriptor.set_properties(ts, props)

		# generated_until represents the earliest time the properties are
		# allowed to change
		generated_until = report_ts + delay + min_duration

		# stop early if enabled for testing 
		if args.end_at and report_ts - start_ts > float(args.end_at):
			break
	
	# content track
	overlay_clip, content_clip = descriptor.generate_content_clip(top_footage, front_footage)
	print(f"length without stills: {util.dur_fmt(content_clip.duration)}")
	overlay_clip, content_clip = add_stills(content_path, content_type, content_fmt, overlay_clip, content_clip)
	print(f"length with stills: {util.dur_fmt(content_clip.duration)}")
	print(f"total generation time: {str(datetime.now() - gen_start)}")

	if not args.yes:
		confirm = input("Render? [y/N] ")
		if confirm != "y":
			print("Exiting")
			exit(0)

	# launch preview application, or save
	if args.preview:
		combined_clip = CompositeVideoClip([content_clip, overlay_clip], size=content_clip.size, use_bgclip=True)
		combined_clip = combined_clip.resize(1)
		combined_clip = combined_clip.subclip(float(args.start_at))
		combined_clip.preview()
	else:
		save(args, overlay_clip, content_clip, content_type)
