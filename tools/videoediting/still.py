import os
import typing

from moviepy.editor import VideoClip, ImageClip, ColorClip, concatenate_videoclips
from moviepy.video.fx import fadein, fadeout

from videoediting.constants import Format
from videoediting.properties.content_property_manager import BasePropertyManager
from videoediting.intro import build_intro
from videoediting.outro import build_outro

import machinepb.machine as pb


def add_stills(
    base_dir: str,
    session_number: int,
    metadata,
    content_type: pb.ContentType,
    property_manager: BasePropertyManager,
    content_plan: pb.ContentTypeStatuses,
    overlay_clip: VideoClip,
    content_clip: VideoClip,
) -> typing.Tuple[VideoClip, VideoClip]:
    FADE_FACTOR = 10

    intro = build_intro(
        base_dir,
        session_number,
        metadata,
        content_type,
        content_plan,
        property_manager.get_stills_config().intro_duration
    )
    intro = intro.fx(fadeout, intro.duration / FADE_FACTOR)

    outro = build_outro(
        base_dir,
        session_number,
        metadata,
        content_type,
        property_manager.get_stills_config().outro_duration
    )
    outro = outro.fx(fadein, outro.duration / FADE_FACTOR)

    overlay_padding = ColorClip(overlay_clip.size, color=(0, 0, 0))

    overlay_clip = concatenate_videoclips([
        overlay_padding.with_duration(intro.duration),
        overlay_clip,
        overlay_padding.with_duration(outro.duration),
    ])
    content_clip = concatenate_videoclips([
        intro,
        content_clip.fadein(intro.duration / FADE_FACTOR).fadeout(outro.duration / FADE_FACTOR),
        outro,
    ])

    return overlay_clip, content_clip
