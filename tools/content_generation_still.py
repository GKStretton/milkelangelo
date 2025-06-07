"""For generating a modified main still image, with branding and session number"""

import argparse
from datetime import datetime
import os
from pathlib import Path

from moviepy.editor import ImageClip, CompositeVideoClip, TextClip
from videoediting.loaders import get_selected_dslr_realpath
from videoediting.compositor_helpers import (
    build_subtitle,
    build_title,
    build_session_number,
)
from videoediting import loaders
from videoediting.constants import Format, MAIN_FONT, PIXEL_FONT


def ts_to_date_time(timestamp):
    # Convert the timestamp to a datetime object in UTC
    dt = datetime.utcfromtimestamp(timestamp)

    # Format the date and time
    formatted = dt.strftime('%Y-%m-%d %H:%M:%S UTC')

    # Return the formatted string
    return formatted


def get_file_date_time(path):
    with open(path + ".creationtime", 'r') as f:
        raw = f.readline()

        # timestamp unix in seconds with decimal
        unix_time = float(raw)

        return ts_to_date_time(unix_time)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("-d", "--base-dir", action="store",
                        help=(
                            "base directory containing session_content and session_metadata"),
                        default="/mnt/md0/light-stores")
    parser.add_argument("-n", "--session-number", action="store",
                        help="session number e.g. 5", required=True)
    args = parser.parse_args()

    base_dir = args.base_dir
    session_number = args.session_number

    metadata = loaders.get_session_metadata(base_dir, session_number)

    duration = 1

    dslr_path = get_selected_dslr_realpath(base_dir, session_number)
    # creationtime is on the raw not post file
    date_time = get_file_date_time(dslr_path.replace("/post/", "/raw/"))
    dslr_clip = ImageClip(dslr_path).with_duration(duration)
    w, h = dslr_clip.size

    align_x, align_y = -1080, -350

    title = TextClip(
        "Milkelangelo",
        size=(1180, 300),
        font_size=120,
        color='white',
        font=MAIN_FONT,
    ).with_duration(duration).with_position((w + align_x, h + align_y))

    date = (
        TextClip(
            date_time,
            size=(1180, 300),
            font_size=70,
            color='white',
            font=MAIN_FONT,
        )
        .with_duration(duration)
        .with_position((w + align_x, h + align_y + 130))
    )

    session_number_text = f"#{metadata['production_id']}" if metadata['production'] else f"dev#{metadata['id']}"
    session_number_clip = TextClip(
        session_number_text,
        size=(800, 200),
        font_size=200,
        color='white',
        font=MAIN_FONT,
        align='west'
    ).with_duration(duration).with_position((30, 50))

    clip = CompositeVideoClip([
        dslr_clip,
        title,
        session_number_clip,
        date
    ], use_bgclip=True)

    # clip.resize((1000, 1000)).show(interactive=True)

    filename = f"{metadata['production_id']}" if metadata['production'] else f"dev{metadata['id']}"
    still_path = os.path.join(base_dir, "session_content", str(
        session_number), "dslr", f"{filename}.jpg")
    clip.save_frame(still_path, t=0.5)

    Path(f"{still_path}.completed").touch()
