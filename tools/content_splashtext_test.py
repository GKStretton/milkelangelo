import numpy as np
import typing
from moviepy.editor import TextClip, CompositeVideoClip, ImageClip, VideoClip
from moviepy.video.fx.resize import resize
from videoediting.constants import Format
from videoediting.loaders import get_session_metadata, get_selected_dslr_image_path
from videoediting.compositor_helpers import build_subtitle, build_title, build_session_number

PIXEL_FONT = "../resources/fonts/MinecraftRegular-Bmg3.otf"
MAIN_FONT = "../resources/fonts/DejaVuSerifCondensed-Italic.ttf"
FONT_SIZE_SUBTITLE = 110
FONT_SIZE_TITLE = 135
FONT_SIZE_SESSION_NUMBER = 170


def pulse(t):
    """
    Function to create a pulse effect.

    1. Sine wave
    2. Take absolute value
    3. 1-curve to make it upside down
    4. trim the bottom off with a max function

    This accents the top of the curve and gives a flat period, creating a beat.
    """
    mag = 0.1
    hz = 2
    scale = 1 + mag * max(1-abs(np.sin(hz * np.pi * t)), 0.2)
    return scale


def slow_grow(t):
    return 1+t/15


def get_size_from_format(fmt: Format) -> typing.Tuple[int, int]:
    if fmt == Format.LANDSCAPE:
        return (1920, 1080)
    elif fmt == Format.PORTRAIT:
        return (1080, 1920)
    else:
        return None


def calculate_splashtext_font_size(text):
    base_size = 80
    base_length = 18  # length of text where font size is at base size

    if len(text) <= base_length:
        return base_size
    else:
        # decrease the font size proportionally to the increase in text length
        return int(base_size * (base_length / len(text)))


def build_splashtext(splash_text, pos, duration) -> VideoClip:
    font_size = calculate_splashtext_font_size(splash_text)
    splash_clip_main = TextClip(splash_text, fontsize=font_size, color='yellow',
                                font=PIXEL_FONT).set_duration(duration)
    splash_clip_shadow = TextClip(splash_text, fontsize=font_size, color='gray',
                                  font=PIXEL_FONT).set_duration(duration).set_position((4, 4))
    final_text = CompositeVideoClip([splash_clip_shadow, splash_clip_main], use_bgclip=True)

    # Apply the pulse effect
    pulsing_clip = resize(final_text, pulse).rotate(20, resample='bicubic')
    w, h = pulsing_clip.size
    x, y = pos

    return pulsing_clip.set_position(lambda t: (x-(w*pulse(t))/2, y-(h*pulse(t))/2))


def build_dslr_image(base_dir: str, session_number: int, duration: float, fmt: Format, pos) -> VideoClip:
    dslr_img = get_selected_dslr_image_path(base_dir, session_number, "selected")
    dslr_clip = ImageClip(dslr_img).set_duration(duration)
    dslr_start_size = min(*get_size_from_format(fmt)) * 0.95
    return (
        dslr_clip
        .fx(resize, (dslr_start_size, dslr_start_size))
        .fx(resize, slow_grow)
        .set_position(pos)
    )


def build_shortform_intro(
    base_dir: str,
    session_number: int,
    metadata,
    fmt: Format,
    duration: float,
    subtitle_text: str,
    splash_text: str = "",
) -> VideoClip:
    title = build_title(('center', 85), duration, font_size=FONT_SIZE_TITLE)

    session_number_clip = build_session_number(
        metadata,
        (20, 230),
        duration,
        font_size=FONT_SIZE_SESSION_NUMBER,
    )

    subtitle = build_subtitle(
        subtitle_text,
        ('center', 1700),
        duration,
        font_size=FONT_SIZE_SUBTITLE
    )

    # Create a composite video clip
    clips = [
        build_dslr_image(base_dir, session_number, duration, fmt, 'center'),
        title,
        session_number_clip,
        subtitle,
    ]

    if splash_text != "":
        pulsing_clip = build_splashtext(splash_text, (700, 320), duration)
        clips.append(pulsing_clip)

    return CompositeVideoClip(clips, size=get_size_from_format(fmt))


if __name__ == "__main__":
    num = 60
    base_dir = "/mnt/md0/light-stores"
    metadata = get_session_metadata(base_dir, num)
    video = build_shortform_intro(
        base_dir,
        num,
        metadata,
        Format.PORTRAIT,
        3.5,
        "Robotic\nArt\nGeneration",
        splash_text="Non-trivial!"
    )
    # video.write_videofile("splash.mp4", fps=60)
    video.resize(0.5).preview()
    # video.resize(0.5).show(interactive=True)
