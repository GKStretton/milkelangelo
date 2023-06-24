import machinepb.machine as pb
from moviepy.editor import ColorClip
from videoediting.constants import Format
from videoediting.compositor_helpers import (
    get_size_from_format
)


def build_outro(
    base_dir: str,
    session_number: int,
    metadata,
    content_type: pb.ContentType,
    duration: int
):
    if content_type == pb.ContentType.CONTENT_TYPE_LONGFORM:
        return ColorClip(get_size_from_format(Format.LANDSCAPE), color=(0, 0, 0)).with_duration(duration)
    if content_type == pb.ContentType.CONTENT_TYPE_SHORTFORM:
        return ColorClip(get_size_from_format(Format.PORTRAIT), color=(0, 0, 0)).with_duration(duration)
    if content_type == pb.ContentType.CONTENT_TYPE_CLEANING:
        return ColorClip(get_size_from_format(Format.PORTRAIT), color=(0, 0, 0)).with_duration(duration)
