"""
This is an automated system for editing the video content of A Study of Light
if this becomes too slow, try using nvenc as moviepy ffmpeg codec
alternately, at using direct [ffmpeg bindings](https://github.com/kkroening/ffmpeg-python)
"""
import argparse
from datetime import datetime
import os

from moviepy.editor import VideoClip, CompositeVideoClip
from pycommon.util import dur_fmt

from videoediting import loaders
from videoediting.footage import FootageWrapper
from videoediting.constants import TOP_CAM, FRONT_CAM
from videoediting.dispense_metadata import DispenseMetadataWrapper
from videoediting.properties.factory import create_property_manager
from videoediting.content_descriptor import ContentDescriptor

import machinepb.machine as pb

# FPS for the final video, not directly related to the source footage
FPS = 60
# how much time to offset the start timestamp of the footage.
# this is basically the latency of the streaming + recording pipeline
# if the lights come on before the state report says so, decrease this magnitude (increase number).
TOP_CAM_TIME_OFFSET = -1.18
FRONT_CAM_TIME_OFFSET = -1.18


def get_args() -> argparse.Namespace:
    """get cli argument"""

    parser = argparse.ArgumentParser()
    parser.add_argument("-d", "--base-dir", action="store",
                        help=("base directory containing session_content and session_metadata"),
                        default="/mnt/md0/light-stores")
    parser.add_argument("-n", "--session-number", action="store",
                        help="session number e.g. 5", required=True)
    parser.add_argument("-e", "--end-at", action="store",
                        help="If set, content will be ended after this timestamp (s)")
    parser.add_argument("-p", "--preview", action="store_true",
                        help="If true, final video will be previewed rather than written")
    parser.add_argument("-x", "--test", action="store_true",
                        help="If true, run test code instead of main functionality")
    parser.add_argument("-t", "--type", action="store",
                        help="content type of output e.g. SHORTFORM | LONGFORM", default="LONGFORM")
    parser.add_argument("-s", "--start-at", action="store",
                        help="set final clip start to this time (s), useful with preview", default=0)
    parser.add_argument("-y", "--yes", action="store_true",
                        help="auto-yes to render confirmation", default=False)

    return parser.parse_args()


def render(
    base_dir: str,
    session_number: int,
    overlay: VideoClip,
    content: VideoClip,
    content_type: pb.ContentType
):
    """ Render overlay, then content, outputting to files """

    output_dir = os.path.join(loaders.get_session_content_path(base_dir, session_number), "video/post/")
    if not os.path.exists(output_dir):
        os.mkdir(output_dir)

    i = 0
    output_file = os.path.join(output_dir, f"{content_type.name}.{i}.mp4")
    while os.path.exists(output_file):
        i += 1
        output_file = os.path.join(output_dir, f"{content_type.name}.{i}.mp4")

    thumbnail_file = os.path.join(output_dir, f"{content_type.name}-thumbnail.{i}.jpg")
    overlay_file = os.path.join(output_dir, f"{content_type.name}-overlay.{i}.mp4")
    content_file = output_file

    # todo: configure thumbnail time
    content.save_frame(thumbnail_file, t=1)
    print(f"wrote thumbnail to {thumbnail_file}")

    overlay_render_start = datetime.now()
    overlay.write_videofile(overlay_file, codec='libx264', fps=FPS)
    print(f"overlay generation time: {str(datetime.now() - overlay_render_start)}")

    content_render_start = datetime.now()
    content.write_videofile(content_file, codec='libx264', fps=FPS)
    print(f"content generation time: {str(datetime.now() - content_render_start)}")


def run():
    gen_start = datetime.now()
    args = get_args()
    print(f"Launching auto_video_post for session {args.session_number} in '{args.base_dir}'\n")

    # the type we're generating for
    content_type = pb.ContentType.from_string(args.type)

    # load data
    session_metadata = loaders.get_session_metadata(args.base_dir, args.session_number)
    state_reports = loaders.get_state_reports(args.base_dir, args.session_number)
    dispense_metadata_wrapper = DispenseMetadataWrapper(args.base_dir, args.session_number)
    misc_data = loaders.get_misc_data(args.base_dir, args.session_number)

    # load webcam footage
    content_path = loaders.get_session_content_path(args.base_dir, args.session_number)
    top_footage = FootageWrapper(
        os.path.join(content_path, "video/raw/" + TOP_CAM),
        timeOffset=TOP_CAM_TIME_OFFSET
    )
    front_footage = FootageWrapper(
        os.path.join(content_path, "video/raw/" + FRONT_CAM),
        timeOffset=FRONT_CAM_TIME_OFFSET
    )

    descriptor = ContentDescriptor(
        session_metadata,
        create_property_manager(content_type),
        dispense_metadata_wrapper,
        misc_data
    )

    # Iterate the state reports, building all the video properties
    descriptor.build_content_descriptor(state_reports, end_at=args.end_at)

    # if it's over the maximium time, do some speed up
    descriptor.limit_duration()

    # generate moviepy clips from the constructed descriptor
    overlay_clip, content_clip = descriptor.generate_content_clip(
        top_footage, front_footage)
    print(f"length without stills: {dur_fmt(content_clip.duration)}")

    # overlay_clip, content_clip = add_stills(
    # content_path, content_type, content_fmt, overlay_clip, content_clip, property_manager)

    print(f"length with stills: {dur_fmt(content_clip.duration)}")
    print(f"total generation time: {str(datetime.now() - gen_start)}")

    # confirm render
    if not args.yes:
        confirm = input("Render? [y/N] ")
        if confirm != "y":
            print("Exiting")
            exit(0)

    # launch preview application, or render
    if args.preview:
        combined_clip = CompositeVideoClip(
            [content_clip, overlay_clip], size=content_clip.size, use_bgclip=True)

        # pylint: disable=E1101
        combined_clip = combined_clip.resize(1)
        combined_clip = combined_clip.subclip(float(args.start_at))
        combined_clip.preview()
    else:
        render(args.base_dir, args.session_number, overlay_clip, content_clip, content_type)


if __name__ == "__main__":
    run()
