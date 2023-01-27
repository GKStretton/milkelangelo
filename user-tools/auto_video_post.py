# This is an automated system for editing the video content of A Study of Light
# if this becomes too slow, try using nvenc as moviepy ffmpeg codec
# alternately, at using direct [ffmpeg bindings](https://github.com/kkroening/ffmpeg-python)
import argparse

from videoediting.content_generator import ContentGenerator
from videoediting.constants import *


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
	print(f"Launching auto_video_post for session {args.session_number} in '{args.base_dir}'\n")

	cg = ContentGenerator(args)

	# cg.generate_content(args, TYPE_LONGFORM)
	cg.generate_content(args, TYPE_SHORTFORM)

	# idea: many types? However, to save on rendering time maybe some of these should be combined.
	# cg.generate_content(LONGFORM_1)
	# cg.generate_content(SHORTFORM_FULL)
	# cg.generate_content(SHORTFORM_HIGHLIGHTS)