from moviepy.editor import VideoClip, CompositeVideoClip
from moviepy.video.fx import crop
import pycommon.image as image

from pycommon.crop_util import CropConfig
import machinepb.machine as pb

from videoediting.properties.content_property_manager import SectionProperties
from videoediting.constants import Scene, Format
from videoediting.compositor_helpers import (
    build_session_number,
    build_speed,
    build_title,
    build_caption
)


def composeLandscape(metadata, props: SectionProperties, top_subclip: VideoClip, front_subclip: VideoClip, caption: str) -> VideoClip:
    landscape_dim = (1920, 1080)

    if props.scene == Scene.DUAL:
        clips = [
            front_subclip.resize(0.7).with_position((50, 'center')),
            top_subclip.resize(1.05).with_position((960, 'center')),

            *build_title((490, 110), top_subclip.duration),
            *build_session_number(metadata, (195, 990), top_subclip.duration),
            build_speed(props.speed, (1700, 20), top_subclip.duration),
        ]
        if caption:
            clips.append(
                build_caption(caption, ('center', 'center'), top_subclip.duration)
            )
        return CompositeVideoClip(clips, size=landscape_dim)

    print("scene {} not supported for landscape format".format(props.scene))
    exit(1)


def composePortrait(metadata, props: SectionProperties, top_subclip: VideoClip, front_subclip: VideoClip, caption: str) -> VideoClip:
    portrait_dim = (1080, 1920)

    if props.scene != Scene.UNDEFINED:
        clips = [
            front_subclip.resize(0.7).with_position(('center', 1180)),
            top_subclip.with_position(('center', 0)),

            *build_title((portrait_dim[0]//2, portrait_dim[1]//2+15), top_subclip.duration, font_size=70),
            *build_session_number(metadata, (10, 0), top_subclip.duration, font_size=80),
            build_speed(props.speed, (900, 10), top_subclip.duration),
        ]
        if caption:
            clips.append(
                build_caption(caption, ('center', 1090), top_subclip.duration)
            )
        return CompositeVideoClip(clips, size=portrait_dim)

    print("scene {} not supported for portrait format".format(props.scene))
    exit(1)


def compositeContentFromFootageSubclips(
    top_subclip: VideoClip,
    top_crop: CropConfig,
    front_subclip: VideoClip,
    front_crop: CropConfig,
    props: SectionProperties,
    fmt: Format,
    session_metadata,
    caption: str
) -> VideoClip:
    #! speed and skip are already applied, this is just for layout!

    # CROP & OVERLAY
    if props.crop:
        if top_crop is not None:
            top_subclip = crop(top_subclip, x1=top_crop.x1, y1=top_crop.y1, x2=top_crop.x2, y2=top_crop.y2)
        if front_crop is not None:
            front_subclip = crop(front_subclip, x1=front_crop.x1, y1=front_crop.y1,
                                 x2=front_crop.x2, y2=front_crop.y2)

        def add_top_overlay(img):
            i = img.copy()
            image.add_overlay(i)
            return i

        def add_front_feather(img):
            i = img.copy()
            image.add_feather(i)
            return i
        if props.vig_overlay:
            top_subclip = top_subclip.image_transform(add_top_overlay)
        if props.front_feather:
            front_subclip = front_subclip.image_transform(add_front_feather)

    # COMPOSITE
    clip: VideoClip = None
    if fmt == Format.LANDSCAPE:
        clip = composeLandscape(session_metadata, props, top_subclip, front_subclip, caption)
    elif fmt == Format.PORTRAIT:
        clip = composePortrait(session_metadata, props, top_subclip, front_subclip, caption)

    return clip
