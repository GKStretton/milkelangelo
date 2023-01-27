from videoediting.constants import *
import pycommon.machine_pb2 as pb

# returns properties for this section. if the second parameter is not 0, 
# this is a "forced_duration". a forced duration requires these properties be
# maintained for this time, even if the state reports change.
def get_section_properties(video_state, state_report, content_type: str) -> dict:
	props = {
		'scene': SCENE_DUAL,
		'speed': 1.0,
		'skip': False,
		'min_duration': 0,
		'format': FORMAT_UNDEFINED,
		'crop': True,
		'overlay': True,
	}

	if content_type == TYPE_LONGFORM:
		props['format'] = FORMAT_LANDSCAPE
	else:
		props['format'] = FORMAT_PORTRAIT


	if state_report.status == pb.WAITING_FOR_DISPENSE:
		props['scene'] = SCENE_DUAL
		# props['skip'] = True
	elif state_report.status == pb.NAVIGATING_IK:
		props['scene'] = SCENE_DUAL
		props['speed'] = 2.5
	elif state_report.status == pb.DISPENSING:
		props['scene'] = SCENE_DUAL
		props['min_duration'] = 3
	else:
		props['scene'] = SCENE_DUAL
		props['speed'] = 10.0
	
	return props
