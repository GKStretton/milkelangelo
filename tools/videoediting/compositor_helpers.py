import math
import numpy as np
import typing

from moviepy.editor import TextClip, ImageSequenceClip, concatenate_videoclips, VideoClip, ImageClip
from moviepy.video.fx.resize import resize
from moviepy.video.fx.loop import loop

from videoediting.constants import Format, MAIN_FONT, PIXEL_FONT
from videoediting.loaders import get_selected_dslr_image_path


def build_title(pos: typing.Tuple[int, int], duration: float, font_size=115) -> TextClip:
    text_size = (1080, 200)
    text_clip = TextClip("A Study of Light", size=text_size, font_size=font_size,
                         color='white', font=MAIN_FONT)

    x = pos[0] if isinstance(pos[0], str) else pos[0] - text_size[0] // 2
    y = pos[1] if isinstance(pos[1], str) else pos[1] - text_size[1] // 2
    text_clip = text_clip.with_position((x, y))
    text_clip = text_clip.with_duration(duration)

    return text_clip


def build_session_number(
        metadata: typing.Dict[str, typing.Any],
        pos: typing.Tuple[typing.Union[int, str],
                          typing.Union[int, str]],
        duration: float,
        center_align: bool = False,
        font_size=115
) -> TextClip:
    text_size = (350, 140)
    session_number_text = f"#{metadata['production_id']}" if metadata['production'] else f"dev#{metadata['id']}"
    text_clip = TextClip(
        session_number_text,
        size=text_size,
        font_size=font_size,
        align='center' if center_align else 'west',
        color='white',
        font=MAIN_FONT
    )

    x = pos[0] if (isinstance(pos[0], str) or not center_align) else pos[0] - text_size[0] // 2
    y = pos[1] if (isinstance(pos[1], str) or not center_align) else pos[1] - text_size[1] // 2
    text_clip = text_clip.with_position((x, y))
    text_clip = text_clip.with_duration(duration)

    return text_clip


def build_subtitle(text: str, pos: typing.Tuple[int, int], duration: float, font_size=80) -> TextClip:
    text_size = (1080, 500)
    text_clip = TextClip(text, size=text_size, font_size=font_size, align='center',
                         color='white', font=MAIN_FONT)

    x = pos[0] if isinstance(pos[0], str) else pos[0] - text_size[0] // 2
    y = pos[1] if isinstance(pos[1], str) else pos[1] - text_size[1] // 2
    text_clip = text_clip.with_position((x, y))
    text_clip = text_clip.with_duration(duration)

    return text_clip


def closest_rounded_speed(speed: float) -> str:
    valid_speeds = [0.5, 1, 1.5, 2, 2.5, 3, 4, 5, 10, 15, 20, 25,
                    30, 40, 50, 75, 100, 150, 200, 300, 400, 500, 750, 1000]
    # minimum difference between speed and one of valid_speeds
    closest_speed = min(valid_speeds, key=lambda x: abs(x - speed))

    if closest_speed == 1 or closest_speed == 2:
        return f"{closest_speed:.0f}x"
    elif closest_speed < 3.0:
        return f"{closest_speed:.1f}x"
    else:
        return f"{closest_speed:.0f}x"


def build_speed(speed: float, pos: typing.Tuple[int, int], duration: float) -> TextClip:
    speed_text = closest_rounded_speed(speed)

    text_clip = TextClip(speed_text, size=(200, 100), font_size=70, align='west',
                         color='white', font=MAIN_FONT)
    text_clip = text_clip.with_position(pos)
    text_clip = text_clip.with_duration(duration)

    return text_clip


def get_size_from_format(fmt: Format) -> typing.Tuple[int, int]:
    if fmt == Format.LANDSCAPE:
        return (1920, 1080)
    elif fmt == Format.PORTRAIT:
        return (1080, 1920)
    else:
        return None


def slow_grow_effect(t):
    """A transition for dslr images, for gradual zooming in"""
    return 1+t/15


def get_drain_effect(duration: float):
    """A transition for dslr images, for drain shrink effect"""
    return lambda t: max(0.01, np.cos(0.5*np.pi*t/duration))


def broom_angle_pulse(t):
    mag = 20
    hz = 2
    return mag * max(1-abs(np.cos(hz * np.pi * t)), 0.2)


def water_resize_pulse(t):
    mag = 0.1
    hz = 2
    return 1+mag * max(1-abs(np.cos(hz * np.pi * t)), 0.2)


def sponge_animation(t):
    mag = 50
    offset = 0.7
    return -np.cos(np.pi * 2 * (t+offset)) * mag


def build_dslr_image(base_dir: str, session_number: int, duration: float, fmt: Format, pos, do_drain_effect: bool = False) -> VideoClip:
    """Return the selected dslr image and slowly zoom in in the clip"""
    dslr_img = get_selected_dslr_image_path(base_dir, session_number, "selected")
    dslr_clip = ImageClip(dslr_img).with_duration(duration)
    dslr_start_size = min(*get_size_from_format(fmt)) * 0.95
    return (
        dslr_clip
        .fx(resize, (dslr_start_size, dslr_start_size))
        .fx(resize, get_drain_effect(duration) if do_drain_effect else slow_grow_effect)
        .with_position(pos)
    )


def build_loader(duration, pos) -> typing.Tuple[VideoClip, VideoClip]:
    """
    returns a bouncer and a countdown, without positions set
    """
    clips = []
    for i in range(math.ceil(duration), 0, -1):
        txt_clip = TextClip(str(i), font_size=200, font=PIXEL_FONT, color='white').with_duration(1)
        clips.append(txt_clip)

    # Concatenate all TextClips into one clip
    countdown = concatenate_videoclips(clips)
    countdown = countdown.subclip(math.ceil(duration) - duration)

    # masking bug with loop
    offset = 0.25
    bounce = ImageSequenceClip("../resources/static_img/loader", fps=60, with_mask=False)
    bounce = loop(bounce, duration=duration+offset).subclip(offset)

    bounce = bounce.with_position(pos)
    countdown = countdown.with_position((pos[0]+67, pos[1]+155))
    return bounce, countdown
