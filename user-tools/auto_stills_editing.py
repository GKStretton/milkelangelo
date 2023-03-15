# For automatically creating still images - thumbnails, outro images, etc.

import argparse
from enum import Enum
from videoediting import loaders
import os
import typing
from videoediting.constants import Format

from PIL import Image, ImageDraw, ImageFont

FONT = "/usr/share/fonts/truetype/dejavu/DejaVuSerifCondensed-Italic.ttf"
FONT_SIZE_TITLE = 130
FONT_SIZE_SUBTITLE = 130
FONT_SIZE_SUBTITLE_OUTRO = 100
FONT_SIZE_SESSION_NUMBER = 200
WHITE = (255, 255, 255)
# SOCIAL ICONS
SI_BASEPATH = "resources/social_icons"
SI_PATHS = [f"{SI_BASEPATH}/youtube.png", f"{SI_BASEPATH}/tiktok.png", f"{SI_BASEPATH}/instagram.png", f"{SI_BASEPATH}/twitter.png"]

TITLE = "A Study of Light"
INTRO_SUBTITLE = "Robotic\nArt\nGeneration"

class StillType(Enum):
	INTRO = 1
	OUTRO = 2

def get_size_from_format(fmt: Format) -> typing.Tuple[int, int]:
	if fmt == Format.LANDSCAPE:
		return (1920, 1080)
	elif fmt == Format.PORTRAIT:
		return (1080, 1920)
	else:
		return None

def get_base_image(metadata, dslr_image: Image.Image, fmt: Format) -> Image.Image:
	size = get_size_from_format(still_format)
	img = Image.new("RGB", size, (0, 0, 0))

	dslr_img_size = min(img.width, img.height)
	dslr_image = dslr_image.resize((dslr_img_size, dslr_img_size))


	dslr_location = (0, 0)
	title_location = (0, 0)
	number_location = (0, 0)
	if fmt == Format.LANDSCAPE:
		dslr_location = (img.width - dslr_image.width, 0)
		title_location = (40, 20)
		number_location = (480, 350)
	if fmt == Format.PORTRAIT:
		dslr_location = (0, (img.height - dslr_image.height) // 2)
		title_location = (60, 20)
		number_location = (img.width / 2, 280)
	
	img.paste(dslr_image, dslr_location)

	draw = ImageDraw.Draw(img)
	# TITLE
	title_font = ImageFont.truetype(FONT, FONT_SIZE_TITLE)
	draw.text(xy=title_location, text=TITLE, fill=WHITE, font=title_font)

	# SESSION NUMBER
	session_number_text = f"#{metadata['production_id']}" if metadata['production'] else f"dev#{metadata['id']}"
	number_font = ImageFont.truetype(FONT, FONT_SIZE_SESSION_NUMBER)
	_, _, w, h = draw.textbbox((0, 0), text=session_number_text, font=number_font)
	draw.text(xy=(number_location[0] - w/2, number_location[1] - h/2), text=session_number_text, fill=WHITE, font=number_font)

	return img

def generate_intro(metadata, dslr_image: Image.Image, still_format: Format) -> Image.Image:
	img = get_base_image(metadata, dslr_image, still_format)

	# Robotic art generation 
	draw = ImageDraw.Draw(img)
	# TITLE
	subtitle_font = ImageFont.truetype(FONT, FONT_SIZE_SUBTITLE)

	subtitle_location = (0, 0)
	if still_format == Format.LANDSCAPE:
		subtitle_location = (490, 765)
	if still_format == Format.PORTRAIT:
		subtitle_location = (img.width / 2, 1675)
	_, _, w, h = draw.textbbox((0, 0), text=INTRO_SUBTITLE, font=subtitle_font)
	draw.text(xy=(subtitle_location[0] - w/2, subtitle_location[1] - h/2), text=INTRO_SUBTITLE, fill=WHITE, font=subtitle_font, align="center")

	return img

def generate_outro(metadata, dslr_image, still_format: Format) -> Image.Image:
	img = get_base_image(metadata, dslr_image, still_format)

	# Robotic art generation 
	draw = ImageDraw.Draw(img)
	# TITLE
	subtitle_font = ImageFont.truetype(FONT, FONT_SIZE_SUBTITLE_OUTRO)

	outro_subtitle = "Follow and\nSubscribe\nFor More!" if still_format == Format.LANDSCAPE else "Follow and Subscribe\nFor More!"
	# Subscribe + follow for more
	subtitle_location = (490, 675) if still_format == Format.LANDSCAPE else (img.width / 2, 1590)
	_, _, w, h = draw.textbbox((0, 0), text=outro_subtitle, font=subtitle_font)
	draw.text(xy=(subtitle_location[0] - w/2, subtitle_location[1] - h/2), text=outro_subtitle, fill=WHITE, font=subtitle_font, align="center")

	# SOCIALS
	si_location = (475, 975) if still_format == Format.LANDSCAPE else (img.width / 2, 1810)
	si_size = 170 if still_format == Format.LANDSCAPE else 160
	si_spacing = 50 if still_format == Format.LANDSCAPE else 70
	
	paste_social_icons(img, si_location, si_size, si_spacing)

	return img

def paste_social_icons(img: Image.Image, loc: typing.Tuple[int, int], size: int, spacing: int):
	# Calculate the number of icons and total width
	num_icons = len(SI_PATHS)
	icons_total_width = size * num_icons + spacing * (num_icons - 1)

	# Calculate the position of the icons
	icons_x = loc[0] - icons_total_width // 2
	icons_y = loc[1] - size // 2

	# Paste the icons onto the thumbnail
	for i in range(num_icons):
		icon_path = SI_PATHS[i]
		icon_image = Image.open(icon_path).resize((size, size))
		icon_x = icons_x + i * (size + spacing)
		icon_y = icons_y
		img.paste(icon_image, (int(icon_x), int(icon_y)))
		

def save_image(base_dir: str, session_number: int, still_type: StillType, still_format: Format, img: Image.Image):
	name = f"{still_type.name}-{still_format.name}.jpg"
	path = os.path.join(base_dir, "session_content", session_number, "stills")
	os.makedirs(path, exist_ok=True)
	img.save(os.path.join(path, name))

if __name__ == "__main__":
	parser = argparse.ArgumentParser()
	parser.add_argument("-d", "--base-dir",
		action="store",
		help="base directory containing session_content and session_metadata",
		required=True
	)
	parser.add_argument("-n", "--session-number",
		action="store",
		help="session number e.g. 5",
		required=True
	)
	parser.add_argument("-i", "--image-choice",
		action="store",
		default="selected",
		help="dslr image name e.g. '205' for 205.jpg"
	)
	parser.add_argument("-t", "--still-type",
		action="store",
		help="Enum number for the StillType",
		required=True
	)
	parser.add_argument("-f", "--still-format",
		action="store",
		help="Enum number for the StillFormat",
		required=True
	)
	parser.add_argument("-p", "--preview",
		action="store_true",
		help="if true, show image rather than save",
	)
	args = parser.parse_args()

	still_type, still_format = StillType.__members__[args.still_type], Format.__members__[args.still_format]

	# load metadata
	metadata = loaders.get_session_metadata(args.base_dir, args.session_number)
	# load dslr image
	dslr_image = loaders.get_selected_dslr_image(args.base_dir, args.session_number, args.image_choice)

	img: Image.Image = None
	if still_type == StillType.INTRO:
		img = generate_intro(metadata, dslr_image, still_format)
	if still_type == StillType.OUTRO:
		img = generate_outro(metadata, dslr_image, still_format)

	if args.preview:
		img.show()
	else:
		save_image(args.base_dir, args.session_number, still_type, still_format, img)