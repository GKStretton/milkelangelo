import moviepy.video.VideoClip as VideoClip
from videoediting.section_properties import SectionProperties
from pycommon.crop_util import CropConfig
from videoediting.constants import *
from moviepy.editor import *
import pycommon.image as image

def composeLandscape(metadata, props: SectionProperties, top_subclip: VideoClip, front_subclip: VideoClip) -> VideoClip:
	landscape_dim = (1920, 1080)

	clip: VideoClip = None
	if props.scene == Scene.DUAL:
		front_subclip = front_subclip.resize(0.7).set_position((50, 'center'))
		top_subclip = top_subclip.resize(1.05).set_position((960, 'center'))
		feather_size = 1
		front_subclip = front_subclip.fx(vfx.fadeout, feather_size, 'north').fx(vfx.fadeout, feather_size, 'south')


		text_clip = TextClip("A Study of Light", size=(900, 140), fontsize=115, color='white', font='DejaVu-Serif-Condensed-Italic')
		text_clip = text_clip.set_position((40, 10))
		text_clip = text_clip.set_duration(top_subclip.duration)

		session_number_text = f"#{metadata['production_id']}" if metadata['production'] else f"dev#{metadata['id']}"
		text_clip2 = TextClip(session_number_text, size=(350, 140), fontsize=115, align='west', color='white', font='DejaVu-Serif-Condensed-Italic')
		text_clip2 = text_clip2.set_position((20, 920))
		text_clip2 = text_clip2.set_duration(top_subclip.duration)

		clip = CompositeVideoClip([front_subclip, top_subclip, text_clip, text_clip2], size=landscape_dim)

	else:
		print("scene {} not supported for landscape format".format(props.scene))
		exit(1)

	return clip

def composePortrait(metadata, props: SectionProperties, top_subclip: VideoClip, front_subclip: VideoClip) -> VideoClip:
	portrait_dim = (1080, 1920)

	clip: VideoClip = None
	if props.scene != Scene.UNDEFINED:
		front_subclip = front_subclip.resize(0.75).set_position(('center', 1120))
		top_subclip = top_subclip.resize(1.05).set_position(('center', 50))

		text_size=(900, 140)
		text_clip = TextClip("A Study of Light", size=text_size, fontsize=115, color='white', font='DejaVu-Serif-Condensed-Italic')
		text_clip = text_clip.set_position(((portrait_dim[0] - text_size[0]) // 2, 955))
		text_clip = text_clip.set_duration(top_subclip.duration)

		session_number_text = "#6854" #f"#{metadata['production_id']}" if metadata['production'] else f"dev#{metadata['id']}"
		text_clip2 = TextClip(session_number_text, size=(350, 120), fontsize=90, align='west', color='white', font='DejaVu-Serif-Condensed-Italic')
		text_clip2 = text_clip2.set_position((20, 10))
		text_clip2 = text_clip2.set_duration(top_subclip.duration)

		clip = CompositeVideoClip([front_subclip, top_subclip, text_clip, text_clip2], size=portrait_dim)
	else:
		print("scene {} not supported for portrait format".format(props.scene))
		exit(1)
	
	return clip

def compositeContentFromFootageSubclips(
		top_subclip: VideoClip,
		top_crop: CropConfig,
		front_subclip: VideoClip,
		front_crop: CropConfig,
		props: SectionProperties,
		fmt: Format,
		session_metadata
	) -> VideoClip:
	#! speed and skip are already applied, this is just for layout!

	# CROP & OVERLAY
	if props.crop:
		if top_crop is not None:
			top_subclip = vfx.crop(top_subclip, x1=top_crop.x1, y1=top_crop.y1, x2=top_crop.x2, y2=top_crop.y2)
		if front_crop is not None:
			front_subclip = vfx.crop(front_subclip, x1=front_crop.x1, y1=front_crop.y1, x2=front_crop.x2, y2=front_crop.y2)
		
		def add_top_overlay(img):
			i = img.copy()
			image.add_overlay(i)
			return i
		def add_front_feather(img):
			i = img.copy()
			image.add_feather(i)
			return i
		if props.vig_overlay:
			top_subclip = top_subclip.fl_image(add_top_overlay)
		if props.front_feather:
			front_subclip = front_subclip.fl_image(add_front_feather)

	# COMPOSITE
	clip: VideoClip = None
	if fmt == Format.LANDSCAPE:
		clip = composeLandscape(session_metadata, props, top_subclip, front_subclip)
	elif fmt == Format.PORTRAIT:
		clip = composePortrait(session_metadata, props, top_subclip, front_subclip)
	
	return clip