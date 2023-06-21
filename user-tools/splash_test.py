import numpy as np
import typing
from moviepy.editor import TextClip, CompositeVideoClip, ImageClip, VideoClip
from moviepy.video.fx.resize import resize
from videoediting.constants import Format
from videoediting.loaders import get_session_metadata, get_selected_dslr_image_path
from videoediting.stills.helpers import get_base_image

pixel_font = "./resources/fonts/MinecraftRegular-Bmg3.otf"


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
    return 1+t/25


def get_size_from_format(fmt: Format) -> typing.Tuple[int, int]:
    if fmt == Format.LANDSCAPE:
        return (1920, 1080)
    elif fmt == Format.PORTRAIT:
        return (1080, 1920)
    else:
        return None


def build_intro(base_dir: str, session_number: int, metadata, fmt: Format, duration: int) -> VideoClip:
    base = np.array(get_base_image(metadata, fmt))
    base_clip = ImageClip(base).set_duration(duration)

    # Define the text
    splash_text = "Milk it for all it's worth!"
    splash_clip_main = TextClip(splash_text, fontsize=70, color='yellow', font=pixel_font).set_duration(duration)
    splash_clip_shadow = TextClip(splash_text, fontsize=70, color='gray',
                                  font=pixel_font).set_duration(duration).set_position((4, 4))
    final_text = CompositeVideoClip([splash_clip_shadow, splash_clip_main], size=(1080, 300))

    # Apply the pulse effect
    pulsing_clip = final_text.resize(pulse).rotate(20, resample='bicubic')
    w, h = pulsing_clip.size
    x, y = 700, 300
    pulsing_clip = pulsing_clip.set_position(lambda t: (x-(w*pulse(t))/2, y-(h*pulse(t))/2))

    # base
    dslr_img = get_selected_dslr_image_path(base_dir, session_number, "selected")
    dslr_clip = ImageClip(dslr_img).set_duration(duration)
    dslr_start_size = min(*get_size_from_format(fmt)) * 0.9
    dslr_clip = dslr_clip.fx(resize, (dslr_start_size, dslr_start_size)).fx(resize, slow_grow)

    # Create a composite video clip
    return CompositeVideoClip([
        base_clip,
        dslr_clip.set_position('center'),
        pulsing_clip
    ], size=get_size_from_format(fmt))


if __name__ == "__main__":
    num = 60
    base_dir = "/mnt/md0/light-stores"
    metadata = get_session_metadata(base_dir, num)
    video = build_intro(base_dir, num, metadata, Format.PORTRAIT, 5)
    video.write_videofile("splash2.mp4", fps=60)
