import moviepy.video.VideoClip as VideoClip
from moviepy.editor import concatenate_videoclips
from videoediting.constants import *
import os
import typing

def add_stills(content_path: str, content_type: ContentType, content_fmt: Format, overlay_clip: VideoClip, content_clip: VideoClip) -> typing.Tuple[VideoClip.VideoClip, VideoClip.VideoClip]:
	path = os.path.join(content_path, "stills")

	introDuration = 1 if content_type == ContentType.SHORTFORM else 3
	introClip = VideoClip.ImageClip(
		img=os.path.join(path, f"INTRO-{content_fmt.name}.jpg"),
		duration=introDuration,
	)
	introClip = introClip.fadeout(introDuration / 5)

	outroDuration = 2 if content_type == ContentType.SHORTFORM else 5
	outroClip = VideoClip.ImageClip(
		img=os.path.join(path, f"OUTRO-{content_fmt.name}.jpg"),
		duration=outroDuration,
	)
	outroClip = outroClip.fadein(outroDuration / 5)

	overlayPadding = VideoClip.ColorClip(overlay_clip.size, color=(0, 0, 0), duration=1)

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