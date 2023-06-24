from videoediting.constants import *
from videoediting.properties.content_property_manager import *
from videoediting.loaders import MiscData


class ShortFormPropertyManager(BasePropertyManager):
    def is_applicable(self, props: SectionProperties) -> bool:
        if props.skip:
            return False
        if props.speed >= 3.0 and props.speed <= 40:
            return True
        return False

    def get_max_content_duration(self) -> typing.Optional[float]:
        return 59.0

    def get_stills_config(self) -> StillsConfig:
        return StillsConfig(
            intro_duration=3.33,
            outro_duration=4,
        )

    def get_format(self) -> Format:
        return Format.PORTRAIT

    def _get_specific_section_properties(
        self,
        current: typing.Tuple[SectionProperties, float, float],
        video_state: VideoState,
        state_report: pb.StateReport,
        dm_wrapper: DispenseMetadataWrapper,
        misc_data: MiscData
    ) -> [SectionProperties, float, float]:
        props, delay, min_duration = current
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
            else:  # dye
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
            props.speed = 15
        elif state_report.pipette_state.dispense_request_number < 1:
            # initial collection and movement is slower
            props.speed = 5
        elif state_report.status == pb.Status.NAVIGATING_IK:
            props.speed = 15
        elif state_report.status == pb.Status.IDLE_STATIONARY:
            props.speed = 100
        else:
            props.speed = 50

        return props, delay, min_duration
