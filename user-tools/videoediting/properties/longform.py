from videoediting.constants import *
from videoediting.properties.content_property_manager import *

class LongFormPropertyManager(BasePropertyManager):
	def get_format(self) -> Format:
		return Format.LANDSCAPE

	def get_section_properties(self, video_state, state_report: pb.StateReport) -> SectionProperties:
		props = self.common_get_section_properties(video_state, state_report)
		if props.skip:
			return props

		# todo: Implement logic specific to LONGFORM ContentType
		return props
