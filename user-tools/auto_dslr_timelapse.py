import argparse
import os
import glob
from moviepy.editor import *
from videoediting.compositor_helpers import *
import videoediting.loaders as loaders

CONTENT_TYPE = "DSLR_TIMELAPSE"
FULL_DIMENSIONS = (1080, 1920)
# dslr footage resize
FOOTAGE_DIMENSIONS = (1080, 1080)
INTERPOLATED_FPS = 60
OUTRO_DURATION = 5

import subprocess

def minterpolate(input_file, output_file):
	command = [
		"ffmpeg",
		"-i", input_file,
		"-vf", f"minterpolate=fps={INTERPOLATED_FPS}:scd=none:mi_mode=mci:me_mode=bidir:me=hexbs:mc_mode=aobmc:vsbmc=1",
		output_file
	]

	subprocess.run(command, check=True)
	print(f"Interpolated video saved to {output_file}")


if __name__ == "__main__":
	parser = argparse.ArgumentParser()
	parser.add_argument("-d", "--base-dir", action="store", help="base directory containing session_content and session_metadata", required=True)
	parser.add_argument("-n", "--session-number", action="store", help="session number e.g. 5", required=True)
	parser.add_argument("-f", "--fps", action="store", type=int, help="dslr frames per second in uninterpolated video, default is 10", default=10)

	args = parser.parse_args()
	print(f"Launching auto_dslr_timelapse.py for session {args.session_number} in '{args.base_dir}'\n")

	session_metadata = loaders.get_session_metadata(args.base_dir, args.session_number)

	path = os.path.join(args.base_dir, "session_content", str(args.session_number), "dslr/post")

	print("Searching for .jpg files...")
	jpg_files = sorted([f for f in glob.glob(os.path.join(path, "*.jpg")) if not f.endswith(".jpg.creationtime") and not os.path.basename(f) == "selected.jpg"])

	print(f"Found {len(jpg_files)} .jpg files. Creating ImageSequenceClip...")
	clip = ImageSequenceClip(jpg_files, fps=args.fps)
	print(f"ImageSequenceClip created with duration {clip.duration:.2f} seconds.")

	print("Compositing video...")
	final_clip = CompositeVideoClip([
		clip.resize(FOOTAGE_DIMENSIONS).set_position(('center', 'center')),

		build_title((FULL_DIMENSIONS[0] // 2, 50), clip.duration),
		build_session_number(session_metadata, (FULL_DIMENSIONS[0] // 2, 220), clip.duration, center=True),
		build_subtitle("Robotic Art Generation\n\nMotion-Interpolated\nDSLR Timelapse", (FULL_DIMENSIONS[0] // 2, 1500), clip.duration),
	], size=FULL_DIMENSIONS)


	output_dir = os.path.join(args.base_dir, "session_content", str(args.session_number), "video/post")

	print("Checking for existing output files...")
	os.makedirs(output_dir, exist_ok=True)

	i = 0
	raw_output_file = os.path.join(output_dir, f"{CONTENT_TYPE}.{i}_raw.mp4")
	while os.path.exists(raw_output_file):
		i += 1
		raw_output_file = os.path.join(output_dir, f"{CONTENT_TYPE}.{i}_raw.mp4")

	print(f"Rendering video to: {raw_output_file}")
	final_clip.write_videofile(raw_output_file, codec='libx264', fps=args.fps)
	print("Video rendering complete.")

	# Run the FFmpeg command to apply the minterpolate filter
	input_file = raw_output_file
	interp_output_file = os.path.join(output_dir, f"{CONTENT_TYPE}.{i}_interpolated{INTERPOLATED_FPS}.mp4")

	print(f"Running FFmpeg to interpolate {input_file}...")
	minterpolate(input_file, interp_output_file)
	print(f"Saved interpolated video to {interp_output_file}.")

	# still concatenation
	print("Concatenating still...")
	interp_video = VideoFileClip(interp_output_file)

	path = os.path.join(args.base_dir, "session_content", str(args.session_number), "stills")
	outroClip = ImageClip(
		img=os.path.join(path, f"OUTRO-PORTRAIT.jpg"),
		duration=OUTRO_DURATION,
	)
	clip_with_outro = concatenate_videoclips([
		interp_video.fadeout(0.2),
		outroClip.fadein(0.5)
	])
	outro_output = os.path.join(output_dir, f"{CONTENT_TYPE}.{i}.mp4")
	clip_with_outro.write_videofile(outro_output, codec='libx264', fps=INTERPOLATED_FPS)
