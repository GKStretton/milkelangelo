from abc import ABC, abstractmethod
from videoediting.section_properties import SectionProperties
from videoediting.constants import *
import machinepb.machine as pb
from videoediting.properties.longform import LongFormPropertyManager
from videoediting.properties.shortform import ShortFormPropertyManager

@dataclass
class SectionProperties:
	scene: Scene = Scene.UNDEFINED
	speed: float = 1.0
	skip: bool = False
	crop: bool = True
	vig_overlay: bool = True
	front_feather: bool = True

	def __str__(self):
		return "\n".join([
			f"{self.scene.name}",
			f"{self.speed}x",
			"skip" if self.skip else "no-skip",
			"crop" if self.crop else "no-crop",
			"vig" if self.vig_overlay else "no-vig",
			"feather" if self.front_feather else "no-feather",
		])

class BasePropertyManager(ABC):

	@abstractmethod
	def get_format(self) -> Format:
		pass

	@abstractmethod
	def get_section_properties(self, video_state, state_report: pb.StateReport) -> SectionProperties:
		pass

	# returns for this section,
	# 1. SectionProperties
	# 2. delay before the properties should come into effect
	# 3. min_duration of these properties
	def common_get_section_properties(self, video_state, state_report: pb.StateReport) -> typing.Tuple[SectionProperties, float, float]:
		props = SectionProperties(
			scene = Scene.DUAL,
			speed = 1.0,
			skip = False,
			crop = True,
			vig_overlay = True,
			front_feather=True,
		)

		if state_report.paused or state_report.status == pb.Status.SLEEPING:
			props.skip = True
			return props

		return props

def create_property_manager(content_type: ContentType) -> BasePropertyManager:
	if content_type == ContentType.LONGFORM:
		return LongFormPropertyManager()
	elif content_type == ContentType.SHORTFORM:
		return ShortFormPropertyManager()
	else:
		raise ValueError(f"Invalid content type: {content_type}")