from videoediting.constants import *
from videoediting.properties.content_property_manager import *
from videoediting.loaders import MiscData

class CleaningPropertyManager(BasePropertyManager):
	def get_max_content_duration(self) -> typing.Optional[float]:
		return 59.0

	def get_stills_config(self) -> StillsConfig:
		return StillsConfig(
			intro_duration=0,
			outro_duration=5,
		)

	def get_format(self) -> Format:
		return Format.PORTRAIT

	def get_section_properties(self, video_state: VideoState, state_report: pb.StateReport, dm_wrapper: DispenseMetadataWrapper, misc_data: MiscData) -> [SectionProperties, float, float]:
		props, delay, min_duration = self.common_get_section_properties(video_state, state_report)
		if props.skip:
			return props, delay, min_duration

		# skip everything before beginning drain
		if video_state.canvas_status != CanvasStatus.AFTER:
			props.skip = True
			return props, delay, min_duration

		if state_report.status == pb.Status.SHUTTING_DOWN:
			props.speed = 2
			return props, delay, min_duration
		
		if (
			state_report.fluid_request.fluid_type == pb.FluidType.FLUID_DRAIN and
			not state_report.fluid_request.complete
		):
			props.max_duration = 30
			props.speed = 2
			return props, delay, min_duration
		
		if (
			state_report.fluid_request.fluid_type == pb.FluidType.FLUID_WATER and
			not state_report.fluid_request.complete
		):
			delay = 3.5
			props.max_duration = 50
			props.speed = 2
			return props, delay, min_duration
		
		if state_report.fluid_request.complete:
			props.skip = True

		return props, delay, min_duration
