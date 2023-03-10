from dataclasses import dataclass
from enum import Enum

TOP_CAM = "top-cam"
FRONT_CAM = "front-cam"

# Scenes for defining composition of top and front cams
class Scene(Enum):
	SCENE_UNDEFINED = 1
	SCENE_FRONT_ONLY = 2
	SCENE_DUAL = 3
	SCENE_TOP_ONLY = 4

class Format(Enum):
	FORMAT_UNDEFINED = 1
	FORMAT_LANDSCAPE = 2
	FORMAT_PORTRAIT = 3

class ContentType(Enum):
	TYPE_LONGFORM = 1
	TYPE_SHORTFORM = 2