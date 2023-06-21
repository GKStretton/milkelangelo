import os
import typing

from moviepy.editor import VideoClip, ImageClip, ColorClip, concatenate_videoclips

from videoediting.constants import Format
from videoediting.properties.content_property_manager import BasePropertyManager

import machinepb.machine as pb


def add_stills(content_path: str, content_type: pb.ContentType, content_fmt: Format, overlay_clip: VideoClip, content_clip: VideoClip, property_manager: BasePropertyManager) -> typing.Tuple[VideoClip, VideoClip]:
    path = os.path.join(content_path, "stills")

    introDuration = property_manager.get_stills_config().intro_duration
    introClip = ImageClip(
        img=os.path.join(path, f"INTRO-{content_fmt.name}.jpg"),
        duration=introDuration,
    )
    # pylint: disable=E1101
    introClip = introClip.fadeout(introDuration / 5)

    outroDuration = property_manager.get_stills_config().outro_duration
    outroClip = ImageClip(
        img=os.path.join(path, f"OUTRO-{content_fmt.name}.jpg"),
        duration=outroDuration,
    )
    # pylint: disable=E1101
    outroClip = outroClip.fadein(outroDuration / 5)

    overlayPadding = ColorClip(overlay_clip.size, color=(0, 0, 0), duration=1)

    overlay_clip = concatenate_videoclips([
        overlayPadding.set_duration(introDuration),
        overlay_clip,
        overlayPadding.set_duration(outroDuration),
    ])
    content_clip = concatenate_videoclips([
        introClip,
        content_clip.fadein(0.5).fadeout(0.5),
        outroClip,
    ])

    return overlay_clip, content_clip
