from videoediting.constants import *
from videoediting.properties.content_property_manager import *
from videoediting.loaders import MiscData

class LongFormPropertyManager(BasePropertyManager):
	def is_applicable(self, props: SectionProperties) -> bool:
		return True

		# todo: Implement logic specific to LONGFORM ContentType
		if props.skip:
			return False
		if props.speed >= 3.0 and props.speed <= 40:
			return True
		return False

	def get_max_content_duration(self) -> typing.Optional[float]:
		# todo: 15 mins?
		return None

	def get_stills_config(self) -> StillsConfig:
		return StillsConfig(
			intro_duration=3,
			outro_duration=5,
		)

	def get_format(self) -> Format:
		return Format.LANDSCAPE

	def get_section_properties(self, video_state: VideoState, state_report: pb.StateReport, dm_wrapper: DispenseMetadataWrapper, misc_data: MiscData) -> typing.Tuple[SectionProperties, float, float]:
		props, delay, min_duration = self.common_get_section_properties(video_state, state_report)
		if props.skip:
			return props, delay, min_duration

		# todo: Implement logic specific to LONGFORM ContentType
		return props, delay, min_duration