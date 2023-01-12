# This is an automated system for editing the video content of A Study of Light
import argparse
import os
from moviepy.editor import *
import yaml

if __name__ == "__main__":
	parser = argparse.ArgumentParser()
	parser.add_argument("-d", "--session-dir", action="store", help="session content directory")
	args = parser.parse_args()

	print(args.session_dir)

	state_reports_path = os.path.join(args.session_dir, "state-reports.yml")
	with open(state_reports_path, 'r') as f:
		state_reports = yaml.load(f, yaml.FullLoader)
	
	print(len(state_reports))
	print(state_reports[0])

	filepath = os.path.join(args.session_dir, "video/raw/top-cam/1.mp4")
	clip = VideoFileClip(filepath)
	clip.preview()