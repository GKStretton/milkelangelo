"""
This is an automated system for editing the video content of A Study of Light
if this becomes too slow, try using nvenc as moviepy ffmpeg codec
alternately, at using direct [ffmpeg bindings](https://github.com/kkroening/ffmpeg-python)
"""
import argparse
import typing
from datetime import datetime
import os
import logging

from moviepy.editor import VideoClip, CompositeVideoClip, clips_array, TextClip
from moviepy.video.fx import resize
from pycommon.util import dur_fmt

from videoediting import loaders
from videoediting.footage import FootageWrapper
from videoediting.constants import TOP_CAM, FRONT_CAM
from videoediting.dispense_metadata import DispenseMetadataWrapper
from videoediting.properties.factory import create_property_manager
from videoediting.content_descriptor import ContentDescriptor
from videoediting.still import add_stills

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
                        help="content type of output e.g. CONTENT_TYPE_SHORTFORM | CONTENT_TYPE_LONGFORM", default="CONTENT_TYPE_LONGFORM")
    parser.add_argument("-s", "--start-at", action="store",
                        help="set final clip start to this time (s), useful with preview", default=0)
    parser.add_argument("-y", "--yes", action="store_true",
                        help="auto-yes to render confirmation")
    parser.add_argument("-f", "--full-duration", action="store_true",
                        help="if true, do not limit video duration to max dur")

    return parser.parse_args()


def setup_logging(base_dir: str, session_number: int, content_type: pb.ContentType):
    _, i = get_output_file(base_dir, session_number, content_type)
    log_file = os.path.join(
        loaders.get_session_content_path(base_dir, session_number),
        "video/post",
        f"{content_type.name}.{i}.log"
    )
    logging.basicConfig(level=logging.INFO,
                        format='%(asctime)s - %(levelname)s - %(message)s',
                        handlers=[
                            logging.StreamHandler(),
                            logging.FileHandler(log_file)
                        ])
    logging.info("logging initialised for %d, %s", session_number, content_type.name)


def get_output_file(base_dir: str, session_number: int, content_type: pb.ContentType) -> typing.Tuple[str, int]:
    """
    Returns main content file and increment number for this content type
    """
    output_dir = os.path.join(loaders.get_session_content_path(base_dir, session_number), "video/post/")
    if not os.path.exists(output_dir):
        os.mkdir(output_dir)

    i = 0
    output_file = os.path.join(output_dir, f"{content_type.name}.{i}.mp4")
    while os.path.exists(output_file):
        i += 1
        output_file = os.path.join(output_dir, f"{content_type.name}.{i}.mp4")

    return output_file, i


def render(
    base_dir: str,
    session_number: int,
    overlay: VideoClip,
    content: VideoClip,
    content_type: pb.ContentType,
    thumbnail_time: int
):
    """ Render overlay, then content, outputting to files """
    output_dir = os.path.join(loaders.get_session_content_path(base_dir, session_number), "video/post/")

    content_file, i = get_output_file(base_dir, session_number, content_type)
    thumbnail_file = os.path.join(output_dir, f"{content_type.name}-thumbnail.{i}.jpg")
    overlay_file = os.path.join(output_dir, f"{content_type.name}-overlay.{i}.mp4")

    content.save_frame(thumbnail_file, t=thumbnail_time, with_mask=False)
    logging.info(f"wrote thumbnail to {thumbnail_file}")

    overlay_render_start = datetime.now()
    overlay.write_videofile(overlay_file, codec='libx264', fps=FPS)
    logging.info(f"overlay generation time: {str(datetime.now() - overlay_render_start)}")

    content_render_start = datetime.now()
    content.write_videofile(content_file, codec='libx264', fps=FPS)
    logging.info(f"content generation time: {str(datetime.now() - content_render_start)}")


def run():
    gen_start = datetime.now()
    args = get_args()
    base_dir, session_number = args.base_dir, int(args.session_number)
    print(f"Launching auto_video_post for session {session_number} in '{base_dir}'\n")

    # the type we're generating for
    content_type = pb.ContentType.from_string(args.type)
    setup_logging(base_dir, session_number, content_type)

    property_manager = create_property_manager(content_type)

    # load data
    session_metadata = loaders.get_session_metadata(base_dir, session_number)
    state_reports = loaders.get_state_reports(base_dir, session_number)
    dispense_metadata_wrapper = DispenseMetadataWrapper(base_dir, session_number)
    misc_data = loaders.get_misc_data(base_dir, session_number)
    content_plan = loaders.get_content_plan(base_dir, session_number)
    profile_snapshot = loaders.get_profile_snapshot(base_dir, session_number)

    # load webcam footage
    content_path = loaders.get_session_content_path(base_dir, session_number)
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
        property_manager,
        dispense_metadata_wrapper,
        misc_data,
        content_plan,
        content_type,
        profile_snapshot
    )

    logging.info("")
    logging.info("** BUILD PROPERTIES LIST (CONTENT DESCRIPTOR)**")
    logging.info("")
    # Iterate the state reports, building all the video properties
    descriptor.build_content_descriptor(state_reports, end_at=args.end_at)

    if not args.full_duration:
        logging.info("")
        logging.info("** LIMIT DURATION **")
        logging.info("")
        # if it's over the maximium time, do some speed up
        descriptor.limit_duration()

    logging.info("")
    logging.info("** GENERATE CLIPS (BUILD SUBCLIPS) **")
    logging.info("")
    # generate moviepy clips from the constructed descriptor
    overlay_clip, content_clip = descriptor.generate_content_clip(
        top_footage, front_footage)

    logging.info("")
    logging.info(f"length without stills: {dur_fmt(content_clip.duration)}")

    overlay_clip, content_clip = add_stills(
        base_dir,
        session_number,
        session_metadata,
        content_type,
        property_manager,
        content_plan,
        overlay_clip,
        content_clip,
    )

    logging.info(f"length with stills: {dur_fmt(content_clip.duration)}")
    logging.info(f"total generation time: {str(datetime.now() - gen_start)}")

    # launch preview application, or render
    if args.preview:
        if args.start_at:
            content_clip = content_clip.subclip(float(args.start_at))
        overlay_clip.fx(resize, 0.5).preview()
    else:
        # confirm render
        if not args.yes:
            confirm = input("Render? [y/N] ")
            if confirm != "y":
                logging.info("Exiting")
                exit(0)

        render(
            base_dir,
            session_number,
            overlay_clip,
            content_clip,
            content_type,
            property_manager.get_stills_config().thumbnail_time
        )


if __name__ == "__main__":
    run()
