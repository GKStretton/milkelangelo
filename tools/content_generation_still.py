"""For generating a modified main still image, with branding and session number"""

import argparse
from datetime import datetime
import os
from pathlib import Path

import numpy as np
from PIL import Image
from moviepy.editor import ImageClip, CompositeVideoClip, TextClip
from videoediting.compositor_helpers import (
    build_subtitle,
    build_title,
    build_session_number,
)
from videoediting import loaders
from videoediting.constants import Format, MAIN_FONT, PIXEL_FONT

# height in pixels of the generated still
OUTPUT_HEIGHT = 2160

# text sizes and margin from the edges, as fractions of the output height.
# Chosen so the text stays in the black corners outside the vignette.
NUMBER_SIZE = 0.07
TITLE_SIZE = 0.035
DATE_SIZE = 0.02
MARGIN = 0.015


def ts_to_date_time(timestamp):
    # Convert the timestamp to a datetime object in UTC
    dt = datetime.utcfromtimestamp(timestamp)

    # Format the date and time
    formatted = dt.strftime('%Y-%m-%d %H:%M:%S UTC')

    # Return the formatted string
    return formatted


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("-d", "--base-dir", action="store",
                        help=(
                            "base directory containing session_content and session_metadata"),
                        default=os.environ.get("LIGHT_STORES_DIR", "/mnt/md0/light-stores"))
    parser.add_argument("-n", "--session-number", action="store",
                        help="session number e.g. 5", required=True)
    args = parser.parse_args()

    base_dir = args.base_dir
    session_number = args.session_number

    metadata = loaders.get_session_metadata(base_dir, session_number)

    duration = 1

    # upscale the frame to a fixed output height, so the text is rendered
    # crisply and the layout looks the same whatever the camera resolution
    still_img = loaders.get_selected_still_image(base_dir, session_number)
    src_h, src_w = still_img.shape[:2]
    h = OUTPUT_HEIGHT
    w = round(src_w * h / src_h)
    still_img = np.array(Image.fromarray(still_img).resize((w, h), Image.LANCZOS))

    date_time = ts_to_date_time(loaders.get_selected_still_timestamp(base_dir, session_number))
    still_clip = ImageClip(still_img).with_duration(duration)

    margin = round(MARGIN * h)

    def text(s: str, size: float) -> TextClip:
        return TextClip(s, font_size=round(size * h), color='white', font=MAIN_FONT).with_duration(duration)

    # session number in the top left
    session_number_text = f"#{metadata['production_id']}" if metadata['production'] else f"dev#{metadata['id']}"
    session_number_clip = text(session_number_text, NUMBER_SIZE).with_position((margin, margin))

    # title with the date under it, right aligned in the bottom right
    title = text("Milkelangelo", TITLE_SIZE)
    date = text(date_time, DATE_SIZE)
    date = date.with_position((w - margin - date.w, h - margin - date.h))
    title = title.with_position((w - margin - title.w, h - margin - date.h - title.h))

    clip = CompositeVideoClip([
        still_clip,
        title,
        session_number_clip,
        date
    ], use_bgclip=True)

    # clip.resize((1000, 1000)).show(interactive=True)

    filename = f"{metadata['production_id']}" if metadata['production'] else f"dev{metadata['id']}"
    out_dir = os.path.join(loaders.get_still_dir(base_dir, session_number), "post")
    os.makedirs(out_dir, exist_ok=True)
    out_path = os.path.join(out_dir, f"{filename}.jpg")
    clip.save_frame(out_path, t=0.5)

    Path(f"{out_path}.completed").touch()
