from videoediting.constants import *
from videoediting.properties.content_property_manager import *
from videoediting.loaders import MiscData

class ShortFormPropertyManager(BasePropertyManager):
	def get_format(self) -> Format:
		return Format.PORTRAIT

	def get_section_properties(self, video_state: VideoState, state_report: pb.StateReport, dm_wrapper: DispenseMetadataWrapper, misc_data: MiscData) -> [SectionProperties, float, float]:
		props, delay, min_duration = self.common_get_section_properties(video_state, state_report)
		if props.skip:
			return props, delay, min_duration

		if (
			video_state.canvas_status == CanvasStatus.AFTER or
			state_report.latest_dslr_file_number >= misc_data.selected_dslr_number
		):
			props.skip = True
			return props, delay, min_duration

		
		# DISPENSE
		if state_report.status == pb.Status.DISPENSING:
			dispense_metadata = dm_wrapper.get_dispense_metadata_from_sr(state_report)
			if dispense_metadata:
				props.skip = dispense_metadata.failed_dispense
				delay = dispense_metadata.dispense_delay_ms / 1000.0
			
			if state_report.pipette_state.vial_held == EMULSIFIER_VIAL:
				min_duration = 2
				props.speed = 1
			else: # dye
				min_duration = 2
				if state_report.pipette_state.dispense_request_number <= 1:
					# first dispense
					props.speed = 1
				else:
					props.speed = 5

			return props, delay, min_duration

		# OTHER
		if state_report.collection_request.request_number < 1:
			props.skip = True
		elif state_report.status == pb.Status.WAITING_FOR_DISPENSE:
			props.speed = 50
		elif state_report.pipette_state.dispense_request_number < 1:
			# initial collection and movement is slower
			props.speed = 4
		elif state_report.status == pb.Status.IDLE_STATIONARY:
			props.speed = 100
		else:
			props.speed = 50

		return props, delay, min_duration
