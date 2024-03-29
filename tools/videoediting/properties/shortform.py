from videoediting.constants import *
from videoediting.properties.content_property_manager import *
from videoediting.loaders import MiscData


class ShortFormPropertyManager(BasePropertyManager):
    def is_applicable(self, props: SectionProperties) -> bool:
        if props.skip:
            return False
        if props.speed >= 3.0:
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
        misc_data: MiscData,
        profile_snapshot: pb.SystemVialConfigurationSnapshot,
        seconds_into_session: float,
    ) -> typing.Tuple[SectionProperties, float, float]:
        props, delay, min_duration = current
        if props.skip:
            return props, delay, min_duration

        if (
            video_state.canvas_status == CanvasStatus.AFTER or
            state_report.latest_dslr_file_number >= misc_data.selected_dslr_number
        ):
            props.skip = True
            return props, delay, min_duration

        # OTHER
        if state_report.collection_request.request_number < 1:
            props.skip = True
            return props, delay, min_duration
        elif state_report.status == pb.Status.WAITING_FOR_DISPENSE:
            props.speed = 10
        elif state_report.pipette_state.dispense_request_number < 2:
            # initial collection and movement is slower
            props.speed = 7
        elif state_report.status == pb.Status.NAVIGATING_IK:
            props.speed = 10
        elif state_report.status == pb.Status.IDLE_STATIONARY:
            base_speed = 50
            cutoff_minutes = 10
            m = seconds_into_session / 60
            if m <= cutoff_minutes:
                props.speed = base_speed
            else:
                props.speed = cutoff_minutes + (m - cutoff_minutes) * 10
        else:
            props.speed = 50

        # DISPENSE
        if state_report.status == pb.Status.DISPENSING:
            dispense_metadata = dm_wrapper.get_dispense_metadata_from_sr(state_report)
            if dispense_metadata:
                props.skip = dispense_metadata.failed_dispense
                delay = dispense_metadata.dispense_delay_ms / 1000.0

            vial_profile = profile_snapshot.profiles.get(str(state_report.pipette_state.vial_held))
            if vial_profile and not vial_profile.footage_ignore:
                min_duration = vial_profile.footage_min_duration_ms / 1000.0
                delay += vial_profile.footage_delay_ms / 1000.0
                props.speed = vial_profile.footage_speed_mult

            # per-dispense overrides
            if dispense_metadata:
                if dispense_metadata.min_duration_override_ms != 0:
                    min_duration = dispense_metadata.min_duration_override_ms / 1000.0
                if dispense_metadata.speed_mult_override != 0:
                    props.speed = dispense_metadata.speed_mult_override

            # override speed of first dispense
            if state_report.pipette_state.dispense_request_number <= 1:
                props.speed = 3
                min_duration = 3

        return props, delay, min_duration
