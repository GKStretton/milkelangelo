from machinepb import machine as pb
from content_plan.loader import *


def buildLongform(base_dir: str, n_str: str) -> pb.ContentTypeStatus:
    ct = pb.ContentType.CONTENT_TYPE_LONGFORM
    raw_title, raw_description = get_random_title_and_description(ct)
    music_file, music_name = get_music(base_dir, ct)

    s = pb.ContentTypeStatus(
        raw_title=raw_title,
        raw_description=raw_description,
        caption=get_caption(ct),
        music_file=music_file,
        music_name=music_name,
    )

    # disabling until it's been shortened
    platform = pb.SocialPlatform.SOCIAL_PLATFORM_YOUTUBE
    s.posts.append(pb.Post(
        platform=platform,
        title=f"{n_str} - {s.raw_title} - Robotic Art Generation Long Cut",
        description=f"{s.raw_description}\n\n{get_hashtags(ct, platform)}\n\n{get_common_text(ct, platform)}\n\n🎵 {music_name}",
        crosspost=False,
        scheduled_unix_timetamp=get_schedule_timestamp(ct),
    ))

    return s
