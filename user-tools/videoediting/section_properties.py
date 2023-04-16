from videoediting.constants import *
import machinepb.machine as pb

@dataclass
class SectionProperties:
	scene: Scene = Scene.UNDEFINED
	speed: float = 1.0
	skip: bool = False
	crop: bool = True
	vig_overlay: bool = True
	front_feather: bool = True

	def __str__(self):
		return f"""{self.scene.name}
{self.speed}x
{"skip" if self.skip else "no-skip"}
{"crop" if self.crop else "no-crop"}
{"vig" if self.vig_overlay else "no-vig"}
{"feather" if self.front_feather else "no-feather"}
"""


# returns for this section,
# 1. SectionProperties
# 2. delay before the properties should come into effect
# 3. min_duration of these properties
def get_section_properties(video_state, state_report: pb.StateReport, content_type: ContentType) -> typing.Tuple[SectionProperties, float, float]:
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

	return props, 0, 0

	# if state_report.status == pb.Status.IDLE_MOVING:
	# 	props.speed = 3.0
	
	# if state_report.status == pb.Status.IDLE_STATIONARY:
	# 	props.speed = 10.0

	#! Do state-based editing once requirements are clearer
	# if state_report.status == pb.WAITING_FOR_DISPENSE:
	# 	props['scene'] = SCENE_DUAL
	# 	# props['skip'] = True
	# elif state_report.status == pb.NAVIGATING_IK:
	# 	props['scene'] = SCENE_DUAL
	# 	props['speed'] = 2.5
	# elif state_report.status == pb.DISPENSING:
	# 	props['scene'] = SCENE_DUAL
	# 	props['min_duration'] = 3
	# else:
	# 	props['scene'] = SCENE_DUAL
	# 	props['speed'] = 10.0
	
