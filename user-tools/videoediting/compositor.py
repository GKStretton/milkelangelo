import moviepy.video.VideoClip as VideoClip
from videoediting.section_properties import SectionProperties
from pycommon.crop_util import CropConfig
from videoediting.constants import *
from moviepy.editor import *
import pycommon.image as image

def compositeContentFromFootageSubclips(
		top_subclip: VideoClip,
		top_crop: CropConfig,
		front_subclip: VideoClip,
		front_crop: CropConfig,
		props: SectionProperties,
		fmt: Format
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
		if props.vig_overlay:
			top_subclip = top_subclip.fl_image(add_top_overlay)

	# COMPOSITE
	landscape_dim = (1920, 1080)
	portrait_dim = (1080, 1920)

	clip: VideoClip = None
	if fmt == Format.LANDSCAPE:
		if props.scene == Scene.DUAL:
			front_subclip = front_subclip.resize(0.65).set_position((10, 'center'))
			top_subclip = top_subclip.resize(1.2).set_position((850, 'center'))
			clip = CompositeVideoClip([front_subclip, top_subclip], size=landscape_dim)
		else:
			print("scene {} not supported".format(props.scene))
			exit(1)
	elif fmt == Format.PORTRAIT:
		if props.scene != Scene.UNDEFINED:
			front_subclip = front_subclip.resize(0.7).set_position(('center', 1150))
			top_subclip = top_subclip.resize(1.2).set_position(('center', 10))
			clip = CompositeVideoClip([front_subclip, top_subclip], size=portrait_dim)
		else:
			print("scene {} not supported".format(props.scene))
			exit(1)
	
	return clip