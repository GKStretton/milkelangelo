from abc import ABC, abstractmethod
from dataclasses import dataclass
import typing
from videoediting.constants import Scene, Format, CanvasStatus
import machinepb.machine as pb
from videoediting.dispense_metadata import DispenseMetadataWrapper
from videoediting.loaders import MiscData

# for keeping track of our state when iterating state reports


@dataclass
class VideoState:
    canvas_status: CanvasStatus = CanvasStatus.BEFORE
    # true if the previous state report was paused
    was_paused: bool = True

    def __str__(self):
        return "\n".join([
            f"CanvasStatus: {self.canvas_status.name}",
        ])


@dataclass
class SectionProperties:
    """
    SectionProperties represents the properties of a section of footage.
    """
    scene: Scene = Scene.UNDEFINED
    speed: float = 1.0
    max_duration: typing.Optional[float] = None
    skip: bool = False
    crop: bool = True
    vig_overlay: bool = True
    front_feather: bool = True

    def __str__(self):
        return "\n".join([
            f"scene: {self.scene.name}",
            f"speed: {self.speed}x",
            f"max_duration: {self.max_duration}",
            "skip" if self.skip else "no-skip",
            "crop" if self.crop else "no-crop",
            "vig" if self.vig_overlay else "no-vig",
            "feather" if self.front_feather else "no-feather",
        ])


@dataclass
class StillsConfig:
    intro_duration: float = 1
    outro_duration: float = 1
    thumbnail_time: float = 1


class BasePropertyManager(ABC):
    """
            Base class for property managers. A property manager is responsible for 
            determining the properties of a section of footage based on:

            - The current state report
            - a video state that is updated by this class to keep track of things
            - other data
    """
    # determines whether a given clip with props may be sped up by the duration limiter

    def is_applicable(self, props: SectionProperties) -> bool:
        return True

    # if set, content will be sped up to fit this duration
    def get_max_content_duration(self) -> typing.Optional[float]:
        return None

    @abstractmethod
    def get_stills_config(self) -> StillsConfig:
        pass

    @abstractmethod
    def get_format(self) -> Format:
        pass

    def get_section_properties(
            self,
            video_state: VideoState,
            state_report: pb.StateReport,
            dm_wrapper: DispenseMetadataWrapper,
            misc_data: MiscData,
            profile_snapshot: pb.SystemVialConfigurationSnapshot
    ) -> typing.Tuple[SectionProperties, float, float]:
        self._update_state_pre(video_state, state_report)


        ts_s = state_report.timestamp_unix_micros / 1.0e6
        seconds_since_session_start = ts_s - misc_data.start_timestamp_s

        cprops, cdelay, cmin_duration = self._common_get_section_properties(video_state, state_report)
        props, delay, min_duration = self._get_specific_section_properties(
            (cprops, cdelay, cmin_duration),
            video_state,
            state_report,
            dm_wrapper,
            misc_data,
            profile_snapshot,
            seconds_since_session_start,
        )
        if props.speed == 0:
            print(ts_s, "PROPS SPEED IS 0")

        self._update_state_post(video_state, state_report)

        return props, delay, min_duration

    # returns for this section,
    # 1. SectionProperties
    # 2. delay before the properties should come into effect
    # 3. min_duration of these properties
    def _common_get_section_properties(
            self,
            video_state: VideoState,
            state_report: pb.StateReport
    ) -> typing.Tuple[SectionProperties, float, float]:
        props = SectionProperties(
            scene=Scene.DUAL,
            speed=1.0,
            skip=False,
            crop=True,
            vig_overlay=True,
            front_feather=True,
        )
        delay, min_duration = 0, 0

        if state_report.paused or state_report.status == pb.Status.SLEEPING:
            props.skip = True
            return props, delay, min_duration

        # delay so that webcam pipelines have started and there is stable footage
        if video_state.was_paused:
            delay = 1

        return props, delay, min_duration

    @abstractmethod
    def _get_specific_section_properties(
            self,
            current: typing.Tuple[SectionProperties, float, float],
            video_state: VideoState,
            state_report: pb.StateReport,
            dm_wrapper: DispenseMetadataWrapper,
            misc_data: MiscData,
            profile_snapshot: pb.SystemVialConfigurationSnapshot
    ) -> SectionProperties:
        """ content-type specific properties logic """

    # for state changes that should happen before this state report has been parsed
    def _update_state_pre(self, video_state: VideoState, state_report: pb.StateReport):
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

    # for state changes that should happen after this state report has been parsed
    def _update_state_post(self, video_state: VideoState, state_report: pb.StateReport):
        video_state.was_paused = state_report.paused
