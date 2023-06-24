import numpy as np
from moviepy.editor import VideoClip, CompositeVideoClip, ImageClip
from moviepy.video.fx import mirror_x, resize, rotate
import machinepb.machine as pb

from videoediting.constants import Format
from videoediting.loaders import get_session_metadata, get_content_plan
from videoediting.splashtext import build_splashtext, splashtext_pulse
from videoediting.compositor_helpers import (
    build_subtitle,
    build_title,
    build_session_number,
    build_loader,
    build_dslr_image,
    get_size_from_format,
    water_resize_pulse,
    sponge_animation,
    slow_grow_effect,
    # broom_angle_pulse,
)

FONT_SIZE_SUBTITLE = 110
FONT_SIZE_TITLE = 135
FONT_SIZE_SESSION_NUMBER = 170

WHIRLPOOL_PATH = "../resources/static_img/whirlpool.png"
BUCKET_RIM_PATH = "../resources/static_img/bucket-rim.png"
BUCKET_MAIN_PATH = "../resources/static_img/bucket-main.png"
BROOM_PATH = "../resources/static_img/broom.png"
SPONGE_PATH = "../resources/static_img/sponge.png"
WATER_PATH = "../resources/static_img/water-emoji.png"
ROBOT_FOREARM_PATH = "../resources/static_img/robot-forearm.png"
ROBOT_UPPERARM_PATH = "../resources/static_img/robot-upperarm.png"


def build_base_intro(metadata, duration, subtitle_text, portrait: bool = False):
    title = build_title(
        ('center', 85) if portrait else (500, 80),
        duration,
        font_size=FONT_SIZE_TITLE if portrait else FONT_SIZE_TITLE - 10
    )

    session_number_clip = build_session_number(
        metadata,
        (20, 230) if portrait else (470, 300),
        duration,
        font_size=FONT_SIZE_SESSION_NUMBER,
        center_align=not portrait
    )

    subtitle = build_subtitle(
        subtitle_text,
        ('center', 1700) if portrait else (470, 870),
        duration,
        font_size=FONT_SIZE_SUBTITLE
    )

    return title, session_number_clip, subtitle


def build_main_intro(
    base_dir: str,
    session_number: int,
    metadata,
    duration: float,
    fmt: Format,
    splash_text: str,
    splash_hue: int,
) -> VideoClip:
    portrait = fmt == Format.PORTRAIT

    title, session_number_clip, subtitle = build_base_intro(
        metadata,
        duration,
        "Robotic\nArt\nGeneration",
        portrait=portrait
    )

    loader_bounce, loader_countdown = build_loader(
        duration,
        (-37, 1600) if portrait else (-40, 770),
    )

    # Create a composite video clip
    clips = [
        build_dslr_image(
            base_dir,
            session_number,
            duration,
            fmt,
            'center' if portrait else lambda t: (870 + (1-slow_grow_effect(t))*500, 'center'),
        ),
        loader_bounce,
        loader_countdown,
        title,
        session_number_clip,
        subtitle,
    ]

    if splash_text != "":
        splash, splash_shadow = build_splashtext(
            splash_text,
            splash_hue,
            (660, 300) if portrait else (470, 515),
            duration)
        clips.append(splash_shadow)
        clips.append(splash)

    def arm_pos(t):
        if portrait:
            return 820, 1520 + np.cos(t*np.pi)*5
        else:
            return 1715, 870 + np.cos(t*np.pi)*5

    arm_scale = 0.5 if portrait else 0.4
    clips.append(
        ImageClip(ROBOT_FOREARM_PATH)
        .fx(resize, arm_scale)
        .with_duration(duration)
        .with_position(arm_pos)
        .fx(rotate, lambda t: np.cos((t+0.5)*np.pi)*16+3, center=(88*arm_scale, 431*arm_scale), expand=False)
    )

    clips.append(
        ImageClip(ROBOT_UPPERARM_PATH)
        .fx(resize, arm_scale)
        .with_duration(duration)
        .with_position(arm_pos)
    )

    return CompositeVideoClip(clips, size=get_size_from_format(fmt))


def build_cleaning_intro(
    base_dir: str,
    session_number: int,
    metadata,
    duration: float,
) -> VideoClip:
    fmt = Format.PORTRAIT

    title, session_number_clip, subtitle = build_base_intro(
        metadata,
        duration,
        "Milk\nArt\nCleanup",
        fmt
    )

    loader_bounce, loader_countdown = build_loader(duration, (910, 1600))

    # Create a composite video clip
    clips = [
        build_dslr_image(base_dir, session_number, duration, fmt, 'center', do_drain_effect=True),
        title,
        session_number_clip,
        subtitle,
        loader_bounce,
        loader_countdown
    ]

    whirlpool_clip = (
        ImageClip(WHIRLPOOL_PATH)
        .with_duration(duration)
        .with_position('center')
    )
    whirlpool_clip = resize(whirlpool_clip, lambda t: max(0.01, np.sin(np.pi*t/duration))/2)
    whirlpool_clip = rotate(whirlpool_clip, lambda t: -t**2*50)
    clips.append(whirlpool_clip)

    # broom = (
    #     ImageClip(BROOM_PATH)
    #     .with_duration(duration)
    # )
    # broom = rotate(broom, broom_angle_pulse, expand=False, center=(450, 50))
    # broom = broom.with_position(lambda t: (800-t*50, 250+t*40))
    # broom = resize(broom, 0.4)
    # clips.append(broom)

    bucket_pos = (0, 1300)
    rim = (
        ImageClip(BUCKET_RIM_PATH)
        .with_duration(duration)
        .with_position(bucket_pos)
    )
    sponge = (
        ImageClip(SPONGE_PATH)
        .fx(resize, 0.8)
        .with_duration(duration)
        .with_position(lambda t: (bucket_pos[0]+130+sponge_animation(t), bucket_pos[1]-30-sponge_animation(t)))
    )
    main = (
        ImageClip(BUCKET_MAIN_PATH)
        .with_duration(duration)
        .with_position(bucket_pos)
    )
    clips.extend([rim, sponge, main])

    clips.append(
        ImageClip(WATER_PATH)
        .with_duration(duration)
        .fx(resize, lambda t: water_resize_pulse(t)*0.6)
        .with_position(lambda t: (800 + (1-water_resize_pulse(t))*300-t*20, 200-(1-water_resize_pulse(t))*300+t*20))
    )

    return CompositeVideoClip(clips, size=get_size_from_format(fmt))


def build_intro(
    base_dir: str,
    session_number: int,
    metadata,
    content_type: pb.ContentType,
    content_plan: pb.ContentTypeStatuses,
    duration: int
):
    splash_text = content_plan.splashtext
    splash_hue = content_plan.splashtext_hue

    if content_type == pb.ContentType.CONTENT_TYPE_LONGFORM:
        return build_main_intro(
            base_dir,
            session_number,
            metadata,
            duration,
            Format.LANDSCAPE,
            splash_text,
            splash_hue
        )
    if content_type == pb.ContentType.CONTENT_TYPE_SHORTFORM:
        return build_main_intro(
            base_dir,
            session_number,
            metadata,
            duration,
            Format.PORTRAIT,
            splash_text,
            splash_hue
        )
    if content_type == pb.ContentType.CONTENT_TYPE_CLEANING:
        return build_cleaning_intro(
            base_dir,
            session_number,
            metadata,
            duration,
        )
