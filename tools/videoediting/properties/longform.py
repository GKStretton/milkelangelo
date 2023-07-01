from videoediting.constants import *
from videoediting.properties.content_property_manager import *
from videoediting.loaders import MiscData


class LongFormPropertyManager(BasePropertyManager):
    def is_applicable(self, props: SectionProperties) -> bool:
        if props.speed > 1.1:
            return True
        return False

    def get_max_content_duration(self) -> typing.Optional[float]:
        return 20*60.0

    def get_stills_config(self) -> StillsConfig:
        return StillsConfig(
            intro_duration=4.33,
            outro_duration=10,
        )

    def get_format(self) -> Format:
        return Format.LANDSCAPE

    def _get_specific_section_properties(
            self,
            current: typing.Tuple[SectionProperties, float, float],
            video_state: VideoState,
            state_report: pb.StateReport,
            dm_wrapper: DispenseMetadataWrapper,
            misc_data: MiscData
    ) -> typing.Tuple[SectionProperties, float, float]:
        props, delay, min_duration = current
        if props.skip:
            return props, delay, min_duration

        # DISPENSE
        if state_report.status == pb.Status.DISPENSING:
            dispense_metadata = dm_wrapper.get_dispense_metadata_from_sr(state_report)
            if dispense_metadata:
                props.skip = dispense_metadata.failed_dispense
                delay = dispense_metadata.dispense_delay_ms / 1000.0

            min_duration = 3

        if state_report.status == pb.Status.WAITING_FOR_DISPENSE:
            # we shouldn't be waiting for dispense. In future could add a timeout
            # that forces dispense after being still a certain time. liquid should
            # be "hot" in this way. also, speed 1 adds to the suspense.
            props.speed = 1

        if (
                state_report.status == pb.Status.IDLE_STATIONARY and
                video_state.canvas_status == CanvasStatus.DURING
        ):
            props.speed = 20

        return props, delay, min_duration
