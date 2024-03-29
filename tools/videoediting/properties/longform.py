import logging
from videoediting.constants import *
from videoediting.properties.content_property_manager import *
from videoediting.loaders import MiscData


class LongFormPropertyManager(BasePropertyManager):
    def is_applicable(self, props: SectionProperties) -> bool:
        if props.speed > 1.1:
            return True
        return False

    def get_max_content_duration(self) -> typing.Optional[float]:
        return 5*60.0

    def get_stills_config(self) -> StillsConfig:
        return StillsConfig(
            intro_duration=4.33,
            outro_duration=15,
        )

    def get_format(self) -> Format:
        return Format.LANDSCAPE

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
        
        if state_report.status == pb.Status.WAKING_UP:
            delay += 1
        
        if not state_report.fluid_request.complete:
            props.speed = 2
        elif state_report.pipette_state.dispense_request_number < 1:
            # initial collection and movement is slower
            props.speed = 1
        elif state_report.status == pb.Status.WAITING_FOR_DISPENSE:
            props.speed = 2
        elif state_report.status == pb.Status.NAVIGATING_IK:
            props.speed = 2
        elif state_report.status == pb.Status.IDLE_STATIONARY:
            base_speed = 20
            cutoff_minutes = 10
            m = seconds_into_session / 60
            if m <= cutoff_minutes:
                props.speed = base_speed
            else:
                props.speed = cutoff_minutes + (m - cutoff_minutes) * 4
        elif state_report.status == pb.Status.IDLE_MOVING:
            props.speed = 5
        else:
            props.speed = 5


        # DISPENSE
        if state_report.status == pb.Status.DISPENSING:
            dispense_metadata = dm_wrapper.get_dispense_metadata_from_sr(state_report)
            if dispense_metadata:
                delay = dispense_metadata.dispense_delay_ms / 1000.0
                # * I think we should include the failures, the longform should be close
                # * to reality, I should fix root problems rather than masking them here.
                # props.skip = dispense_metadata.failed_dispense

            # we use min_duration to prevent early speedup when system goes idle
            vial_profile = profile_snapshot.profiles.get(str(state_report.pipette_state.vial_held))
            if vial_profile and not vial_profile.footage_ignore:
                props.speed = 1
                min_duration = vial_profile.footage_min_duration_ms / 1000.0

            # per-dispense min_duration override
            if dispense_metadata:
                if dispense_metadata.min_duration_override_ms != 0:
                    min_duration = dispense_metadata.min_duration_override_ms / 1000.0

            # Basically forces speed 1 for at least 30 seconds after the latest dispense
            min_duration = min(20, min_duration)

        return props, delay, min_duration
