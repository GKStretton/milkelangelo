from dataclasses import dataclass
from enum import Enum

TOP_CAM = "top-cam"
FRONT_CAM = "front-cam"

EMULSIFIER_VIAL = 4

# Scenes for defining composition of top and front cams
class Scene(Enum):
	UNDEFINED = 1
	FRONT_ONLY = 2
	DUAL = 3
	TOP_ONLY = 4

class Format(Enum):
	UNDEFINED = 1
	LANDSCAPE = 2
	PORTRAIT = 3

class ContentType(Enum):
	LONGFORM = 1
	SHORTFORM = 2