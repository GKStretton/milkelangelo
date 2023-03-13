# For automatically creating still images - thumbnails, outro images, etc.

import argparse
from enum import Enum
from videoediting import loaders
import os
import typing
from videoediting.constants import Format

from PIL import Image

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

def get_base_image(fmt: Format) -> Image.Image:
	size = get_size_from_format(still_format)
	return Image.new("RGB", size, (0, 0, 0))

def generate_intro(metadata, dslr_image, still_format: Format) -> Image.Image:
	img = get_base_image(still_format)
	# todo: image logic
	return img

def generate_outro(metadata, dslr_image, still_format: Format) -> Image.Image:
	img = get_base_image(still_format)
	# todo: image logic
	return img

def save_image(base_dir: str, session_number: int, still_type: StillType, still_format: Format, img: Image.Image):
	name = f"{still_type.name}-{still_format.name}.jpg"
	path = os.path.join(base_dir, "session_content", session_number, "stills", name)
	img.save(path)

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
	parser.add_argument("t", "--still-type",
		action="store",
		help="Enum number for the StillType",
		required=True
	)
	parser.add_argument("f", "--still-format",
		action="store",
		help="Enum number for the StillFormat",
		required=True
	)
	args = parser.parse_args()

	still_type, still_format = StillType(args.still_type), Format(args.still_format)

	# load metadata
	metadata = loaders.get_session_metadata(args.base_dir, args.session_number)
	# load dslr image
	dslr_image = loaders.get_selected_dslr_image(args.base_dir, args.session_number)

	img = None
	if still_type == StillType.INTRO:
		img = generate_intro(metadata, dslr_image)
	if still_type == StillType.OUTRO:
		img = generate_outro(metadata, dslr_image)

	save_image(args.base_dir, args.session_number, still_type, still_format, img)