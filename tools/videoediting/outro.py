import numpy as np
import machinepb.machine as pb
from moviepy.editor import ColorClip, CompositeVideoClip, TextClip, ImageClip, ImageSequenceClip
from moviepy.video.fx import resize, rotate, loop
from videoediting.constants import Format, MAIN_FONT, PIXEL_FONT
from videoediting.compositor_helpers import (
    get_size_from_format,
    build_dslr_image,
    slow_grow_effect,
    build_title,
    build_session_number,
)

SHOW_TWITCH_BANNER = True


def timezones():
    from datetime import datetime, time
    from dateutil.relativedelta import relativedelta
    import pytz

    # Create a datetime object for start time UTC
    dt_utc = datetime.combine(datetime.today() + relativedelta(weeks=1), time(17, 0)).replace(tzinfo=pytz.UTC)

    # Convert to UK time
    zones = [
        'America/Los_Angeles',
        'America/New_York',
        # 'America/Sao_Paulo',
        'UTC',
        'Europe/London',
        'Europe/Berlin',
        # 'Asia/Shanghai',
        # 'Asia/Kolkata',
        # 'Australia/Sydney',
        # 'Africa/Johannesburg',
        # 'Asia/Tokyo',
    ]

    outputs = []
    for z in zones:
        dt = dt_utc.astimezone(pytz.timezone(z))

        # Format the output
        time = dt.strftime('%H:%M')
        zone = dt.strftime('%Z')
        if zone == '-03':
            zone = "BRT"
        output = (time, zone)
        print(output)
        outputs.append(output)

    return outputs


FONT_SIZE_TITLE = 135
FONT_SIZE_SESSION_NUMBER = 80


def build_twitch_block(pos, duration):
    clips = []
    x, y = pos[0], pos[1]

    clips.append(
        ImageClip("../resources/static_img/twitch-block.png")
        .with_duration(duration)
        .with_position(pos)
    )

    tl = ImageClip("../resources/social_icons/twitch.png")
    tl_scale = 0.5
    tl_size = int(tl.size[0] / 2 * tl_scale)
    clips.append(
        tl
        .fx(resize, tl_scale)
        .with_duration(duration)
        .with_position((x+20, y+80))
        .fx(rotate, lambda t: np.cos(t*np.pi)*10, expand=False, center=(tl_size, tl_size))
    )

    clips.append(
        ImageSequenceClip("../resources/static_img/stream_anim", fps=60)
        .fx(loop, duration=duration)
        .with_position((x+750, y+50))
    )

    clips.append(
        TextClip(
            "Join LIVE every Sunday",
            # color=twitch_color,
            font_size=70,
            font=PIXEL_FONT
        )
        .with_position((x+160, y+10))
        .with_duration(duration)
    )
    clips.append(
        TextClip(
            "twitch.tv/StudyOfLight",
            # color=twitch_color,
            font_size=70,
            font=PIXEL_FONT
        )
        .with_position((x+195, y+125))
        .with_duration(duration)
    )

    for i, t in enumerate(timezones()):
        print(i, t)
        clips.append(
            TextClip(
                f"{t[1]}\n{t[0]}",
                # color=twitch_color,
                font_size=60,
                font=PIXEL_FONT
            )
            .with_position(lambda t, i=i: (x+430 + (i-2)*140, y+230 + np.cos(np.pi * t + 0.4*(-i)*np.pi) * 10))
            .with_duration(duration)
        )

    return clips


def build_outro(
    base_dir: str,
    session_number: int,
    metadata,
    content_type: pb.ContentType,
    duration: int
):
    fmt = Format.PORTRAIT
    if content_type == pb.ContentType.CONTENT_TYPE_LONGFORM:
        fmt = Format.LANDSCAPE
    portrait = fmt == Format.PORTRAIT

    clips = []

    clips.append(
        build_dslr_image(
            base_dir,
            session_number,
            duration,
            fmt,
            lambda t: ('center', 300 + (1-slow_grow_effect(t))*500)
            if fmt == Format.PORTRAIT else (870 + (1-slow_grow_effect(t))*500, 'center')
        )
    )

    clips.extend(
        build_title(
            (540, 85) if portrait else (500, 80),
            duration,
            font_size=FONT_SIZE_TITLE if portrait else FONT_SIZE_TITLE - 10
        )
    )

    clips.extend(
        build_session_number(
            metadata,
            (540, 230) if portrait else (500, 210),
            duration,
            font_size=FONT_SIZE_SESSION_NUMBER,
            center_align=True
        )
    )

    if SHOW_TWITCH_BANNER:
        clips.extend(
            build_twitch_block(
                (60, 1250) if portrait else (10, 390),
                duration
            )
        )

    # SOCIAL ICONS
    SI_BASEPATH = "../resources/social_icons"
    SI_PATHS = [
        f"{SI_BASEPATH}/empty.png",
        f"{SI_BASEPATH}/youtube.png",
        # f"{SI_BASEPATH}/tiktok.png",
        f"{SI_BASEPATH}/instagram.png",
        # f"{SI_BASEPATH}/twitter.png",
        f"{SI_BASEPATH}/empty.png",
    ]

    for i, p in enumerate(SI_PATHS):
        clips.append(
            ImageClip(p)
            .fx(resize, 0.5)
            .with_duration(duration)
            .with_position((100+i*230, 1750) if portrait else (65+i*230, 890))
        )

    clips.extend(
        draw_main_text(
            "Follow for more!",
            (250, 1630) if portrait else (200, 760),
            80,
            duration
        )
    )

    clips.extend(
        draw_main_text(
            "Thank you for watching!",
            (100, 280) if portrait else (50, 270),
            80,
            duration
        )
    )

    return CompositeVideoClip(clips, size=get_size_from_format(fmt))


def draw_main_text(text, pos, font_size, duration):
    offset = 3
    shadow = (
        TextClip(text, color="black", font=MAIN_FONT, font_size=font_size)
        .with_duration(duration)
        .with_position((pos[0]+offset, pos[1]+offset))
    )
    main_text = (
        TextClip(text, color="white", font=MAIN_FONT, font_size=font_size)
        .with_duration(duration)
        .with_position(pos)
    )
    return [shadow, main_text]
