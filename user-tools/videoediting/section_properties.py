from videoediting.constants import *
import pycommon.machine_pb2 as pb

@dataclass
class SectionProperties:
	scene: Scene = Scene.SCENE_UNDEFINED
	speed: float = 1.0
	skip: bool = False
	fmt: Format = Format.FORMAT_UNDEFINED
	crop: bool = True
	vig_overlay: bool = True

	def __str__(self):
		return f"""{self.fmt}
{self.scene.name}
{self.speed}x
{"skip" if self.skip else "no-skip"}
{"crop" if self.crop else "no-crop"}
{"vig" if self.vig_overlay else "no-vig"}
"""


# returns properties for this section. if the second parameter is not 0, 
# this is a "forced_duration". a forced duration requires these properties be
# maintained for this time, even if the state reports change.
def get_section_properties(video_state, state_report, content_type: ContentType) -> SectionProperties:
	props = SectionProperties(
		scene = Scene.SCENE_DUAL,
		speed = 1.0,
		skip = False,
		fmt = Format.FORMAT_UNDEFINED,
		crop = True,
		vig_overlay = True,
	)

	if content_type == ContentType.TYPE_LONGFORM:
		props.fmt = Format.FORMAT_LANDSCAPE
	else:
		props.fmt = Format.FORMAT_PORTRAIT

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
	
	return props
