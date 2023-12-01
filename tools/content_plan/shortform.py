from machinepb import machine as pb
from content_plan.loader import *


def buildShortform(n_str: str) -> pb.ContentTypeStatus:
    ct = pb.ContentType.CONTENT_TYPE_SHORTFORM
    raw_title, raw_description = get_random_title_and_description(ct)

    s = pb.ContentTypeStatus(
        raw_title=raw_title,
        raw_description=raw_description,
        caption=get_caption(ct),
    )

    platform = pb.SocialPlatform.SOCIAL_PLATFORM_YOUTUBE
    s.posts.append(pb.Post(
        platform=platform,
        title=append_title_hashtags(f"{s.raw_title} - {n_str}", ct, platform),
        description=f"{get_hashtags(ct, platform)}\n\n{s.raw_description}\n\n{get_common_text(ct, platform)}",
        crosspost=False,
        scheduled_unix_timetamp=get_schedule_timestamp(ct),
    ))

    # platform = pb.SocialPlatform.SOCIAL_PLATFORM_TIKTOK
    # s.posts.append(pb.Post(
    #     platform=platform,
    #     title=append_title_hashtags(f"{s.raw_title} - {n_str}", ct, platform) + "\n\n" + get_common_text(ct, platform),
    #     description="N/A",
    #     crosspost=False,
    #     scheduled_unix_timetamp=get_schedule_timestamp(ct),
    # ))

    platform = pb.SocialPlatform.SOCIAL_PLATFORM_INSTAGRAM
    s.posts.append(pb.Post(
        platform=platform,
        title=f"{s.raw_title} - {n_str}\n\n{s.raw_description}\n\n{get_common_text(ct, platform)}\n\n{get_hashtags(ct, platform)}",
        description="N/A",
        crosspost=False,
        scheduled_unix_timetamp=get_schedule_timestamp(ct),
    ))

    # platform = pb.SocialPlatform.SOCIAL_PLATFORM_FACEBOOK
    # s.posts.append(pb.Post(
    #     platform=platform,
    #     title=f"{s.raw_title} - {n_str}\n\n{s.raw_description}\n\n{get_common_text(ct, platform)}\n\n{get_hashtags(ct, platform)}",
    #     description="N/A",
    #     crosspost=False,
    #     scheduled_unix_timetamp=get_schedule_timestamp(ct),
    # ))

    # platform = pb.SocialPlatform.SOCIAL_PLATFORM_TWITTER
    # s.posts.append(pb.Post(
    #     platform=platform,
    #     title=f"{s.raw_title} - {n_str}\n\n{get_common_text(ct, platform)}\n\n{get_hashtags(ct, platform)}",
    #     description="N/A",
    #     crosspost=False,
    #     scheduled_unix_timetamp=get_schedule_timestamp(ct),
    # ))

    return s
