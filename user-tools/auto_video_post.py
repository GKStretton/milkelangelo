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
from termcolor import colored
from google.protobuf.json_format import MessageToJson



TOP_CAM = "top-cam"
FRONT_CAM = "front-cam"

def get_session_metadata(args: argparse.Namespace):
	filename = "{}_session.yml".format(args.session_number)
	path = os.path.join(args.base_dir, "session_metadata", filename)
	yml = None
	with open(path, 'r') as f:
		yml = yaml.load(f, Loader=yaml.FullLoader)
	print("Loaded session metadata: {}\n".format(yml))
	return yml
	
def get_session_content_path(args: argparse.Namespace):
	return os.path.join(args.base_dir, "session_content", args.session_number)

def get_state_reports(args: argparse.Namespace):
	content_path = get_session_content_path(args)
	state_reports = None
	with open(os.path.join(content_path, "state-reports.yml"), 'r') as f:
		state_reports = yaml.load(f, yaml.FullLoader)
	print("Loaded {} state report entries\n".format(len(state_reports)))
	return state_reports

# Scenes for defining composition of top and front cams
SCENE_UNDEFINED = "UNDEFINED"
SCENE_FRONT_ONLY = "FRONT_ONLY"
SCENE_FRONT_PRIMARY = "FRONT_PRIMARY"
SCENE_TOP_PRIMARY = "TOP_PRIMARY"
SCENE_TOP_ONLY = "TOP_ONLY"

FORMAT_UNDEFINED = "UNDEFINED"
FORMAT_LANDSCAPE = "LANDSCAPE"
FORMAT_PORTRAIT = "PORTRAIT"

TYPE_LONGFORM = "LONGFORM"
TYPE_SHORTFORM = "SHORTFORM"

class ContentGenerator:
	def __init__(self, args: argparse.Namespace):
		self.args = args

		## Gather resources
		print("~~~ GATHERING RESOURCES ~~~")

		self.session_metadata = get_session_metadata(args)
		self.state_reports = get_state_reports(args)

		if not args.dry_run:
			# load camera footage
			content_path = get_session_content_path(args)
			
			self.top_footage = FootageWrapper(os.path.join(content_path, "video/raw/" + TOP_CAM))
			self.front_footage = FootageWrapper(os.path.join(content_path, "video/raw/" + FRONT_CAM))

		print()

	def generate_content(self, args: argparse.Namespace, content_type: str):
		print("Iterating state reports...")
		subclips: typing.List[VideoClip.VideoClip] = []

		# state for the content generation, used with get_section_properties
		video_state = {}

		# content already generated up to this point, so skip any state reports without info past this point
		section_start_ts = 0
		section_properties = {
			'scene': SCENE_UNDEFINED,
			'speed': 1.0,
			'skip': False,
		}
		previous_state_report = None

		sections = 0

		for i in range(len(self.state_reports)):
			# get status report object
			report = ParseDict(self.state_reports[i], pb.StateReport())

			# get status name and ts
			status_name = pb.Status.Name(report.status)
			report_ts = float(report.timestamp_unix_micros) / 1.0e6

			# i    timestamp_str     ts       STATUS_X
			print("{}\t{}     ({})\t{}".format(colored(str(i), attrs=['bold', 'underline']), util.ts_format(report_ts), colored("{:.2f}".format(report_ts), 'red'), status_name))

			new_section_properties = self.get_section_properties(video_state, report, content_type)
			if new_section_properties != section_properties or (args.end_at is not None and int(args.end_at) <= i): # if content properties have changed, or we're exiting early
				print("\t*Property change*")
				if i == 0 or section_properties['skip']:
					print("\t{} content until {}".format(colored("Skipping", attrs=['bold']), colored("{:.2f}".format(report_ts), 'red')))
					print()
				elif section_start_ts >= report_ts:
					# Skip this report if we have already generated up to it. This may occur after a force_duration.
					print("\t{} content for section because section_start_ts {:.2f} is already greater than this report timestamp".format(colored("Skipping", attrs=['bold']), section_start_ts))
				else:
					print("\t{} content up to this SR: {} -> {}\t({:.2f})".format(colored("Generating", attrs=['bold']), colored("{:.2f}".format(section_start_ts), 'green'), colored("{:.2f}".format(report_ts), 'red'), report_ts-section_start_ts))
					print()
					sections += 1
					# todo: add properties
					if not args.dry_run:
						clip = self.generate_section(section_start_ts, report_ts, section_properties, previous_state_report)
						if clip is not None:
							subclips.append(clip)

				# update state
				section_start_ts = max(section_start_ts, report_ts)
				section_properties = new_section_properties
				previous_state_report = report
				print("\tUpdated section_start_ts to {} and properties to {}".format(colored("{:.2f}".format(section_start_ts), 'green'), section_properties))
			else:
				# print("\tSkipping state report because video properties have not changed")
				pass

			# force duration of this section if min_duration is set
			min_duration = new_section_properties['min_duration']
			if min_duration > 0:
				print("Forcing min_duration from {:.2f} to {:.2f} ({:.2f})".format(section_start_ts, section_start_ts + min_duration, min_duration))
				if not args.dry_run:
					# todo: add properties
					clip = self.generate_section(section_start_ts, section_start_ts + min_duration, section_properties, previous_state_report)
					if clip is not None:
						subclips.append(clip)

				section_start_ts += min_duration
			
			if args.end_at is not None and int(args.end_at) <= i:
				print("i exceeds args.end_at, breaking loop")
				break


		# Final section handling

		if section_properties['skip'] or args.end_at is not None:
			print("{}\tSkipping final section from {} to {}".format(colored("end", attrs=['bold', 'underline']), colored("{:.2f}".format(section_start_ts), 'green'), colored("end_of_footage", 'red')))
		else:
			print("{}\tGenerating final section from {} to {}".format(colored("end", attrs=['bold', 'underline']), colored("{:.2f}".format(section_start_ts), 'green'), colored("end_of_footage", 'red')))
			sections += 1

			# todo: add properties
			if not args.dry_run:
				clip = self.generate_section(section_start_ts, None, section_properties, previous_state_report)
				if clip is not None:
					subclips.append(clip)
				print()
		
		# Summary
		
		print("-"*40)
		n_sr = len(self.state_reports) if args.end_at is None else int(args.end_at) + 1
		print("State reports: {}\nSections: {}".format(n_sr, sections))
		print("-"*40)
		print()

		# Writing

		if not args.dry_run:
			print("Concatenating...")
			final_clip = concatenate_videoclips(subclips)

			if args.preview:
				final_clip.preview()
			elif args.show:
				final_clip.show(interactive=True)
			else:
				self.write_final_clip(content_type, final_clip)
	
	# returns properties for this section. if the second parameter is not 0, 
	# this is a "forced_duration". a forced duration requires these properties be
	# maintained for this time, even if the state reports change.
	def get_section_properties(self, video_state, state_report, content_type: str) -> typing.Tuple[dict, float]:
		props = {
			'scene': SCENE_FRONT_PRIMARY,
			'speed': 1.0,
			'skip': False,
			'min_duration': 0,
			'format': FORMAT_UNDEFINED
		}

		if content_type == TYPE_LONGFORM:
			props['format'] = FORMAT_LANDSCAPE
		else:
			props['format'] = FORMAT_PORTRAIT

		if state_report.status == pb.WAITING_FOR_DISPENSE:
			props['skip'] = True
		elif state_report.status == pb.NAVIGATING_IK:
			props['scene'] = SCENE_TOP_PRIMARY
			props['speed'] = 2.5
		elif state_report.status == pb.DISPENSING:
			props['scene'] = SCENE_TOP_PRIMARY
			props['min_duration'] = 3
		else:
			props['scene'] = SCENE_FRONT_PRIMARY
			props['speed'] = 5.0
		
		return props
	

	# generates a section of the content, one subclip.
	def generate_section(self, start_ts: float, end_ts: float, properties: dict, state_report) -> VideoClip:
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

		speed = properties['speed']
		scene = properties['scene']
		format = properties['format']

		# SPEED
		if speed != 1:
			top_clip = top_clip.speedx(speed)
			front_clip = front_clip.speedx(speed)
		
		landscape_dim = (1920, 1080)
		portrait_dim = (1080, 1920)

		clip = None

		if format == FORMAT_LANDSCAPE:
			if scene == SCENE_FRONT_PRIMARY:
				top_clip = top_clip.resize(0.5).set_position((50, 50))
				clip = CompositeVideoClip([front_clip, top_clip], size=landscape_dim)
			elif scene == SCENE_TOP_PRIMARY:
				front_clip = front_clip.resize(0.5).set_position((50, 50))
				clip = CompositeVideoClip([top_clip, front_clip], size=landscape_dim)
			else:
				print("scene {} not supported".format(scene))
				return None
		elif format == FORMAT_PORTRAIT:
			if scene == SCENE_FRONT_PRIMARY or scene == SCENE_TOP_PRIMARY:
				top_clip = top_clip.resize(0.5).set_position((100, 100))
				front_clip = front_clip.resize(0.5).set_position((100, 1000))
				clip = CompositeVideoClip([front_clip, top_clip], size=portrait_dim)
			else:
				print("scene {} not supported".format(scene))
				return None
		else:
			print("format {} not supported, returning".format(format))
			return None
		
		if clip is None:
			return None
		
		if self.args.inspect:
			# add text to clip with debug
			#! commenting out this line in /etc/ImageMagick-6/policy.xml was required:
			#! <policy domain="path" rights="none" pattern="@*" />
			#! https://github.com/Zulko/moviepy/issues/401#issuecomment-278679961

			msg = MessageToJson(
				state_report,
				including_default_value_fields=True,
			)

			text = "\n".join([
				"PROPERTIES",
				"  speed:        {:.2f}x".format(speed),
				"  scene:        {}".format(scene),
				"  format:       {}".format(format),
				"  min_duration: {}".format(properties['min_duration']),
				"STATE REPORT",
				msg
			])

			txt: TextClip = TextClip(text, font='DejaVu-Sans-Mono-Bold', fontsize=15, color='white', align='West')
			txt = txt.set_position((5, 5)).set_duration(clip.duration)
			clip = CompositeVideoClip([clip, txt])

		return clip

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
	parser.add_argument("-x", "--dry-run", action="store_true", help="If true, this will be a dry run and no content will be generated")
	parser.add_argument("-e", "--end-at", action="store", help="If set, content will be ended at this state report")
	parser.add_argument("-p", "--preview", action="store_true", help="If true, final video will be previewed rather than written")
	parser.add_argument("-s", "--show", action="store_true", help="If true, show start of final clip")

	args = parser.parse_args()
	print("Launching auto_video_post for session {} in '{}'\n".format(args.session_number, args.base_dir))

	cg = ContentGenerator(args)

	# cg.test()

	cg.generate_content(args, TYPE_LONGFORM)
	# cg.generate_content(LONGFORM_1)
	# cg.generate_content(SHORTFORM_FULL)
	# cg.generate_content(SHORTFORM_HIGHLIGHTS)