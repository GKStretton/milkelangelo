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
        return None
        # return 20*60.0

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
        profile_snapshot: pb.SystemVialConfigurationSnapshot
    ) -> typing.Tuple[SectionProperties, float, float]:
        props, delay, min_duration = current
        if props.skip:
            return props, delay, min_duration

        if state_report.status == pb.Status.WAKING_UP:
            delay += 1

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
                min_duration = vial_profile.footage_min_duration_ms / 1000.0

            # per-dispense min_duration override
            if dispense_metadata:
                if dispense_metadata.min_duration_override_ms != 0:
                    min_duration = dispense_metadata.min_duration_override_ms / 1000.0

            # Basically forces speed 1 for at least 30 seconds after the latest dispense
            min_duration = min(20, min_duration)

        if state_report.status == pb.Status.WAITING_FOR_DISPENSE:
            # we shouldn't be waiting for dispense. In future could add a timeout
            # that forces dispense after being still a certain time. liquid should
            # be "hot" in this way. also, speed 1 adds to the suspense.
            props.speed = 1

        if (
                state_report.status == pb.Status.IDLE_STATIONARY and
                state_report.fluid_request.complete
        ):
            props.speed = 20

        return props, delay, min_duration
