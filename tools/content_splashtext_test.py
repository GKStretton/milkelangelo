import random
import typing
import numpy as np
from PIL import Image

from moviepy.editor import TextClip, CompositeVideoClip, ImageClip, VideoClip, ImageSequenceClip
from moviepy.video.fx.resize import resize
from moviepy.video.fx.loop import loop
from videoediting.constants import Format
from videoediting.loaders import get_session_metadata, get_selected_dslr_image_path
from videoediting.compositor_helpers import build_subtitle, build_title, build_session_number

# PIXEL_FONT = "../resources/fonts/MinecraftRegular-Bmg3.otf"
PIXEL_FONT = "../resources/fonts/monogram.ttf"
CS_FONT = "../resources/fonts/Comic-Sans-MS.ttf"
MAIN_FONT = "../resources/fonts/DejaVuSerifCondensed-Italic.ttf"
FONT_SIZE_SUBTITLE = 110
FONT_SIZE_TITLE = 135
FONT_SIZE_SESSION_NUMBER = 170
FONT_RESCALE = 3


def pulse(t):
    """
    Function to create a pulse effect.

    1. Cosine wave
    2. Take absolute value
    3. 1-curve to make it upside down
    4. trim the bottom off with a max function

    This accents the top of the curve and gives a flat period, creating a beat.
    """
    mag = 0.1
    hz = 2
    scale = 1 + mag * max(1-abs(np.cos(hz * np.pi * t)), 0.2)
    return scale / FONT_RESCALE


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
    base_size = 80 if get_splashfont(text) == CS_FONT else 100
    base_length = 20  # length of text where font size is at base size

    if len(text) <= base_length:
        return base_size
    else:
        # decrease the font size proportionally to the increase in text length
        return int(base_size * (base_length / len(text)))


def generate_hue():
    """
    generate a hue for the splashtext, excluding yellow and dark blue
    """
    choices = []
    for i in range(0, 360):
        if i >= 50 and i <= 80:
            continue
        if i >= 225 and i <= 260:
            continue
        choices.append(i)
    return random.choice(choices)


def get_splashfont(text: str) -> str:
    cs_texts = ["Now in your favorite font!", "Dare to be different!"]
    if text in cs_texts:
        return CS_FONT
    return PIXEL_FONT


def build_splashtext(splash_text, pos, duration) -> typing.Tuple[VideoClip, VideoClip]:
    """returns main text and shadow"""
    angle = -15
    hue = generate_hue()
    color = f"hsv({hue}, 255, 255)"
    shadow_color = f"hsv({hue}, 50, 40)"
    print(color)

    font_size = calculate_splashtext_font_size(splash_text)
    splash_clip_main = TextClip(splash_text, font_size=font_size*FONT_RESCALE, color=color,
                                font=get_splashfont(splash_text)).with_duration(duration)
    splash_clip_shadow = TextClip(splash_text, font_size=font_size*FONT_RESCALE, color=shadow_color,
                                  font=get_splashfont(splash_text)).with_duration(duration)

    splash_clip_main = splash_clip_main.rotate(angle, resample='bilinear')
    splash_clip_main = resize(splash_clip_main, pulse)

    w, h = splash_clip_main.size
    x, y = pos
    splash_clip_main = splash_clip_main.with_position(lambda t: (
        x-(w*pulse(t)*FONT_RESCALE)/2, y-(h*pulse(t)*FONT_RESCALE)/2))

    # Apply the pulse effect to shadow
    xoffset = -4
    yoffset = 3
    pulsing_shadow = resize(splash_clip_shadow, pulse).rotate(angle, resample='bicubic')
    w, h = pulsing_shadow.size
    x2, y2 = pos[0]+xoffset, pos[1]+yoffset
    pulsing_shadow = pulsing_shadow.with_position(lambda t: (
        x2-(w*pulse(t)*FONT_RESCALE)/2, y2-(h*pulse(t)*FONT_RESCALE)/2))

    return splash_clip_main, pulsing_shadow


def build_dslr_image(base_dir: str, session_number: int, duration: float, fmt: Format, pos) -> VideoClip:
    dslr_img = get_selected_dslr_image_path(base_dir, session_number, "selected")
    dslr_clip = ImageClip(dslr_img).with_duration(duration)
    dslr_start_size = min(*get_size_from_format(fmt)) * 0.95
    return (
        dslr_clip
        .fx(resize, (dslr_start_size, dslr_start_size))
        .fx(resize, slow_grow)
        .with_position(pos)
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

    offset = 0.25
    # masking bug with loop
    loader = ImageSequenceClip("../resources/static_img/loader", fps=60, with_mask=False)
    loader = loop(loader, duration=duration+offset).subclip(offset).with_position((0, 1570))

    # Create a composite video clip
    clips = [
        build_dslr_image(base_dir, session_number, duration, fmt, 'center'),
        title,
        session_number_clip,
        subtitle,
        loader,
    ]

    if splash_text != "":
        splash, splash_shadow = build_splashtext(splash_text, (660, 300), duration)
        clips.append(splash_shadow)
        clips.append(splash)
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
        splash_text="Dancing is not science!"
    )
    # video.write_videofile("splash2.mp4", fps=60)
    video.resize(1).preview()
    # video.resize(0.5).show(interactive=True)
