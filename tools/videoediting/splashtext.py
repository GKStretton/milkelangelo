import random
import typing

import numpy as np
from moviepy.editor import VideoClip, TextClip
from moviepy.video.fx.resize import resize

from videoediting.constants import PIXEL_FONT, CS_FONT

FONT_RESCALE = 3


def splashtext_pulse(t):
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


def calculate_splashtext_font_size(text):
    base_size = 80 if get_splashfont(text) == CS_FONT else 100
    base_length = 20  # length of text where font size is at base size

    if len(text) <= base_length:
        return base_size
    else:
        # decrease the font size proportionally to the increase in text length
        return int(base_size * (base_length / len(text)))


def get_splashfont(text: str) -> str:
    cs_texts = ["Now in your favorite font!", "Dare to be different!"]
    if text in cs_texts:
        return CS_FONT
    return PIXEL_FONT


def build_splashtext(splash_text: str, hue: int, pos, duration) -> typing.Tuple[VideoClip, VideoClip]:
    """returns main text and shadow"""
    angle = -15
    color = f"hsv({hue}, 255, 255)"
    shadow_color = f"hsv({hue}, 50, 40)"
    print(color)

    font_size = calculate_splashtext_font_size(splash_text)
    splash_clip_main = TextClip(splash_text, font_size=font_size*FONT_RESCALE, color=color,
                                font=get_splashfont(splash_text)).with_duration(duration)
    splash_clip_shadow = TextClip(splash_text, font_size=font_size*FONT_RESCALE, color=shadow_color,
                                  font=get_splashfont(splash_text)).with_duration(duration)

    splash_clip_main = splash_clip_main.rotate(angle, resample='bilinear')
    splash_clip_main = resize(splash_clip_main, splashtext_pulse)

    w, h = splash_clip_main.size
    x, y = pos
    splash_clip_main = splash_clip_main.with_position(lambda t: (
        x-(w*splashtext_pulse(t)*FONT_RESCALE)/2, y-(h*splashtext_pulse(t)*FONT_RESCALE)/2))

    # Apply the pulse effect to shadow
    xoffset = 3
    yoffset = 3
    pulsing_shadow = resize(splash_clip_shadow, splashtext_pulse).rotate(angle, resample='bicubic')
    w, h = pulsing_shadow.size
    x2, y2 = pos[0]+xoffset, pos[1]+yoffset
    pulsing_shadow = pulsing_shadow.with_position(lambda t: (
        x2-(w*splashtext_pulse(t)*FONT_RESCALE)/2, y2-(h*splashtext_pulse(t)*FONT_RESCALE)/2))

    return splash_clip_main, pulsing_shadow
