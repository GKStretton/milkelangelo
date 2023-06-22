from moviepy.editor import TextClip, CompositeVideoClip
import numpy as np


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
    scale = 1 + mag * max(1-abs(np.cos(hz * np.pi * t)), 0.2)
    return scale


PIXEL_FONT = "../resources/fonts/MinecraftRegular-Bmg3.otf"
CS_FONT = "../resources/fonts/Comic-Sans-MS.ttf"
MAIN_FONT = "../resources/fonts/DejaVuSerifCondensed-Italic.ttf"
duration = 5
font_size = 80

splash_text = "What door?"
color = 'purple'

splash_clip_main = (
    TextClip(splash_text, font_size=font_size, color=color, font=CS_FONT)
    .with_duration(duration)
    .resize(pulse)
    .rotate(-15)
)

CompositeVideoClip([splash_clip_main], use_bgclip=False).preview()
