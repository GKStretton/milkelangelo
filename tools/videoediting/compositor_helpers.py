from moviepy.video.VideoClip import TextClip
import typing

MAIN_FONT = "../resources/fonts/DejaVuSerifCondensed-Italic.ttf"


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
