import typing
import dataclasses
import json

from moviepy.editor import VideoClip, TextClip, concatenate_videoclips, clips_array, CompositeVideoClip

from betterproto import Casing
import machinepb.machine as pb

from videoediting.properties.content_property_manager import VideoState, SectionProperties, BasePropertyManager
from videoediting.compositor import compositeContentFromFootageSubclips
from videoediting.footage import FootageWrapper
from videoediting.dispense_metadata import DispenseMetadataWrapper
from videoediting.loaders import MiscData

from pycommon.util import ts_fmt, ts_format, floats_are_equal

FILTER_UNCHANGED = True
FINAL_DURATION = 2


class ContentDescriptor:
    """
        This is a descriptor (list of timestamps and video properties) for a single
        piece of content (video). It is used to generate the moviepy aggregate
        clips for the overlay and content.
    """

    def __init__(
        self,
        session_metadata,
        property_manager: BasePropertyManager,
        dispense_metadata_wrapper: DispenseMetadataWrapper,
        misc_data: MiscData,
        content_plan: pb.ContentTypeStatuses,
        content_type: pb.ContentType,
        profile_snapshot: pb.SystemVialConfigurationSnapshot
    ):
        self.session_metadata = session_metadata
        self.property_manager = property_manager
        self.dispense_metadata_wrapper = dispense_metadata_wrapper
        self.misc_data = misc_data
        self.content_plan = content_plan
        self.content_type = content_type
        self.caption = self.content_plan.content_statuses[self.content_type.name].caption
        self.profile_snapshot = profile_snapshot

        # [(timestamp_s, sr), ...]
        self.state_reports: typing.List[typing.Tuple[float, pb.StateReport, VideoState]] = []
        # [(timestamp_s, props), ...]
        self.properties: typing.List[typing.Tuple[float, SectionProperties]] = []

    def set_state_report(self, timestamp: float, state_report: pb.StateReport, video_state: VideoState):
        if len(self.state_reports) > 0 and self.state_reports[-1][0] > timestamp:
            print(
                f"set_state_report with timestamp {timestamp}, but previously seen state report has timestamp {self.state_reports[-1][0]}")
            exit(1)

        self.state_reports.append((timestamp, state_report, dataclasses.replace(video_state)))

    def set_properties(self, timestamp: float, properties: SectionProperties):
        if len(self.properties) > 0 and self.properties[-1][0] > timestamp:
            print(
                f"set_properties with timestamp {timestamp}, but previously seen state report has timestamp {self.properties[-1][0]}")
            exit(1)

        # only add prop if it's different to last
        if FILTER_UNCHANGED and len(self.properties) > 0 and self.properties[-1][1] == properties:
            return

        self.properties.append((timestamp, properties))

    def _generate_state_report_overlay_clip(self) -> VideoClip:
        """
        This only looks at state reports, and generates a full uncut clip of all
        state reports as text. No skipping or speed changes take place here.

        This is done to simplify the main generation code which iterates properties.
        Because state reports can occur when there is no property change, only
        iterating properties doesn't let you change the state reports in the overlay
        at the right times. You could loop both state reports and properties and
        generate an overlay that way, but I couldn't work out how to do that, this
        way is easier but less efficient.
        """
        if len(self.properties) == 0 or len(self.state_reports) == 0:
            print("_generate_raw_overlay_clip found no properties / state_reports")
            exit(1)

        start_timestamp = self.state_reports[0][0]

        print(f"gen raw state report overlay for {len(self.state_reports)} state reports")
        # generate raw state report overlay in normal time
        sr_clips = []
        for i, _ in enumerate(self.state_reports):
            print(i)
            timestamp, sr, video_state = self.state_reports[i]

            # show last one for this long
            duration = FINAL_DURATION
            if i < len(self.state_reports) - 1:
                next_timestamp, _, _ = self.state_reports[i+1]
                duration = next_timestamp - timestamp

            sr_fmt = json.dumps(
                sr.to_dict(include_default_values=True, casing=Casing.SNAKE),
                indent=4,
                sort_keys=True,
            )
            text_str = "STATE REPORT:\n"+ts_format(timestamp) + "\n" + sr_fmt
            text_str += "\n\nVIDEO STATE:\n" + str(video_state) + "\n"
            txt: TextClip = TextClip(text_str, font='DejaVu-Sans-Mono', font_size=10, color='white', align='West')
            txt = txt.with_duration(duration)

            sr_clips.append(txt)

        sr_full_clip = concatenate_videoclips(sr_clips)
        return sr_full_clip

    def generate_content_clip(self, top_footage: FootageWrapper, front_footage: FootageWrapper) -> typing.Tuple[VideoClip, VideoClip]:
        """
        self.properties has been generated, so we have the list
        [
            (timestamp, Props)
        ]

        So we need to iterate that list (self.properties) and create a subclip
        for each property. only max_duration is enforced here (because that's
        simple truncation), the other functionality like min_duration is handled
        in the property generation stage.


        """
        if len(self.properties) == 0:
            print("generate_content_clip found no properties / state_reports")
            exit(1)

        # Generate full raw overlay from the properties
        raw_overlay = self._generate_state_report_overlay_clip()
        start_timestamp = self.properties[0][0]

        print(f"gen content for {len(self.properties)} props")
        overlay_clips = []
        content_clips = []
        for i, _ in enumerate(self.properties):
            print()
            print(i)
            props = self.properties[i][1]
            print(props)
            if props.skip:
                print("skipping")
                continue

            ts1_abs = self.properties[i][0]
            ts1_rel = ts1_abs - start_timestamp
            # default to end
            ts2_abs, ts2_rel = None, None
            # if there's another prop after this
            if i < len(self.properties) - 1:
                ts2_abs = self.properties[i+1][0]
                ts2_rel = ts2_abs - start_timestamp

            # enforce max_duration by limiting ts2
            if self.properties[i][1].max_duration is not None:
                ts2_abs = min(ts2_abs, ts1_abs + self.properties[i][1].max_duration)
                ts2_rel = ts2_abs - start_timestamp

            # build overlay subclip
            overlay_raw_subclip: VideoClip = raw_overlay.subclip(ts1_rel, ts2_rel)
            text_str = "PROPS:\n"+ts_format(self.properties[i][0]) + "\n" + str(props)
            # commenting out this line in /etc/ImageMagick-6/policy.xml was required:
            # <policy domain="path" rights="none" pattern="@*" />
            # https://github.com/Zulko/moviepy/issues/401#issuecomment-278679961
            txt: TextClip = TextClip(text_str, font='DejaVu-Sans-Mono', font_size=10, color='white', align='West')
            txt = txt.with_duration(overlay_raw_subclip.duration)
            if props.speed != 1.0:
                overlay_raw_subclip = overlay_raw_subclip.multiply_speed(props.speed)
                txt = txt.multiply_speed(props.speed)
            txt = txt.with_position((0, 600))
            overlay_subclip = CompositeVideoClip([overlay_raw_subclip, txt], size=(400, 800))

            # build footage subclips
            print(
                f"Getting top_footage between {ts1_abs} and {ts2_abs} ({ts_fmt(ts1_rel)} to {ts_fmt(ts2_rel)})")
            top_subclip, top_crop = top_footage.get_subclip(ts1_abs, ts2_abs)
            print(
                f"Getting front_footage between {ts1_abs} and {ts2_abs} ({ts_fmt(ts1_rel)} to {ts_fmt(ts2_rel)})")
            front_subclip, front_crop = front_footage.get_subclip(ts1_abs, ts2_abs)

            # apply speed to both
            if props.speed != 1.0:
                top_subclip = top_subclip.multiply_speed(props.speed)
                front_subclip = front_subclip.multiply_speed(props.speed)

            # clips should all be same length unless it's the last property.
            if i != len(self.properties) - 1 and not floats_are_equal(0.00001, [overlay_subclip.duration, top_subclip.duration, front_subclip.duration]):
                # these should be same length if the footage has been padded correctly''
                print("processed subclips are not same duration: {} {} {}, exiting".format(
                    overlay_subclip.duration, top_subclip.duration, front_subclip.duration))
                exit(1)

            content_clips.append(compositeContentFromFootageSubclips(top_subclip, top_crop,
                                 front_subclip, front_crop, props, self.property_manager.get_format(), self.session_metadata, self.caption))
            overlay_clips.append(overlay_subclip)

        if len(content_clips) == 0:
            print("no content clips built, skipped all?")
            exit(1)

        return concatenate_videoclips(overlay_clips), concatenate_videoclips(content_clips)

    def limit_duration(self):
        """
        Adjusts the speed of applicable sections in a video to meet a maximum duration constraint.

        This function calculates the total duration of a video, including stills, and compares it 
        to the maximum allowed duration specified in the property manager. If the total duration 
        exceeds the maximum allowed duration, the function increases the playback speed of 
        applicable sections (based on the `is_applicable` criteria) to reduce the total duration.

        The function skips sections marked with the 'skip' property and only considers sections 
        with speeds between 3.0 and 40.0 for adjustment.

        Parameters:
        - property_manager (BasePropertyManager): An object containing properties and configurations,
        including the maximum allowed duration for the video.

        Returns:
        None
        """

        max_duration = self.property_manager.get_max_content_duration()
        if max_duration is None:
            print("max_duration pref is none")
            return

        # calculate existing durations
        total_duration = 0
        applicable_duration = 0
        for i in range(len(self.properties) - 1):
            ts, props = self.properties[i]
            ts_next, props_next = self.properties[i+1]

            if not props.skip:
                maximum = props.max_duration if props.max_duration is not None else 1e10
                duration = min(maximum, ts_next - ts) / props.speed
                total_duration += duration
                if self.property_manager.is_applicable(props):
                    applicable_duration += duration

        stills_time = self.property_manager.get_stills_config().intro_duration + self.property_manager.get_stills_config().outro_duration
        full_total = total_duration + stills_time

        reduction = full_total - max_duration
        threshold = 0.1
        if reduction <= threshold:
            print(f"No duration reduction needed: {full_total:.2f}s <= {max_duration + threshold:.2f}s")
            return

        new_applicable_time = applicable_duration - reduction
        # Find speed factor to achieve the reduction
        speed_factor = applicable_duration / new_applicable_time

        print(f"Target reduction of {reduction:.2f}s, {speed_factor:.2f} factor per applicable property")
        print(f"Calculated total duration: {full_total:.2f}s, applicable: {applicable_duration:.2f}s.")

        # apply reduction
        for _, v in enumerate(self.properties):
            _, props = v
            if props.skip or not self.property_manager.is_applicable(props):
                continue
            props.speed *= speed_factor

    def build_content_descriptor(
        self,
        state_reports: typing.Tuple[float, pb.StateReport],
        end_at: typing.Optional[float]
    ):
        """
        This iterates the state reports to build a list of properties.
        Depending on the state, the video properties (speed, skip toggle) will
        change. The properties are fetched from the property_manager.

        The non-trivial element here is the handling of "delay" and "min_duration"
        functionalities. These result in property transitions occuring when
        there is no state report transition. For example:
            - delay of 500ms on a new speed means preserve the old speed until
              500ms after the state change.
            - min_duration of 5s means that the properties generated for this
              state must be maintained for at least 5s, even if an upcoming
              state report suggests changing props before then.

        Min duration is handled with some precise logic around a
        "generated_until" tracker, which serves as a block on property changes
        until that time. Delay is handled alongside that mechanism.
        """
        start_ts = state_reports[0][0]

        # no properties may be set until after this time
        generated_until = 0

        # unused
        video_state = VideoState()

        # Iterate state reports to build descriptor
        for i, (report_ts, report) in enumerate(state_reports):
            # Get section properties
            props, delay, min_duration = self.property_manager.get_section_properties(
                video_state,
                report,
                self.dispense_metadata_wrapper,
                self.misc_data,
                self.profile_snapshot
            )

            # Always set state report in descriptor for use in the overlays
            self.set_state_report(report_ts, report, video_state)

            # if we've generated beyond everything this report covers, skip it
            if (
                i+1 < len(state_reports) and
                generated_until >= state_reports[i+1][0] and
                generated_until >= report_ts + delay + min_duration
            ):
                continue

            # only change properties after generated_until
            ts = max(generated_until, report_ts + delay)
            self.set_properties(ts, props)

            # generated_until represents the earliest time the properties are
            # allowed to change
            generated_until = report_ts + delay + min_duration

            # stop early if enabled for testing
            if end_at and report_ts - start_ts > float(end_at):
                break
