from abc import ABC, abstractmethod
import typing
from videoediting.constants import *
import machinepb.machine as pb
from videoediting.dispense_metadata import DispenseMetadataWrapper
from videoediting.loaders import MiscData

# for keeping track of our state when iterating state reports
@dataclass
class VideoState:
	canvas_status: CanvasStatus = CanvasStatus.BEFORE

	def __str__(self):
		return "\n".join([
			f"CanvasStatus: {self.canvas_status.name}",
		])

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

@dataclass
class StillsConfig:
	intro_duration: float = 1
	outro_duration: float = 1

class BasePropertyManager(ABC):
	# if set, content will be sped up to fit this duration
	def get_max_content_duration(self) -> typing.Optional[float]:
		return None

	@abstractmethod
	def get_stills_config(self) -> StillsConfig:
		pass

	@abstractmethod
	def get_format(self) -> Format:
		pass

	@abstractmethod
	def get_section_properties(self, video_state: VideoState, state_report: pb.StateReport, dm_wrapper: DispenseMetadataWrapper, misc_data: MiscData) -> SectionProperties:
		pass

	# returns for this section,
	# 1. SectionProperties
	# 2. delay before the properties should come into effect
	# 3. min_duration of these properties
	def common_get_section_properties(self, video_state: VideoState, state_report: pb.StateReport) -> typing.Tuple[SectionProperties, float, float]:
		self.update_state(video_state, state_report)

		props = SectionProperties(
			scene = Scene.DUAL,
			speed = 1.0,
			skip = False,
			crop = True,
			vig_overlay = True,
			front_feather=True,
		)
		delay, min_duration = 0, 0

		if state_report.paused or state_report.status == pb.Status.SLEEPING:
			props.skip = True
			return props, delay, min_duration


		return props, delay, min_duration

	def update_state(self, video_state: VideoState, state_report: pb.StateReport):
		# canvas status
		if video_state.canvas_status == CanvasStatus.BEFORE:
			if (
				state_report.fluid_request.fluid_type == pb.FluidType.FLUID_MILK and
				state_report.fluid_request.complete and
				not state_report.fluid_request.open_drain
			):
				video_state.canvas_status = CanvasStatus.DURING
		elif video_state.canvas_status == CanvasStatus.DURING:
			if (
				state_report.fluid_request.fluid_type != pb.FluidType.FLUID_MILK or
				state_report.fluid_request.open_drain
			):
				video_state.canvas_status = CanvasStatus.AFTER
