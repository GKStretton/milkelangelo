from machinepb import machine as pb
from content_plan.loader import *


def buildStill(n_str: str) -> pb.ContentTypeStatus:
    ct = pb.ContentType.CONTENT_TYPE_STILL

    s = pb.ContentTypeStatus(
        raw_title=f"Piece {n_str} - Milkelangelo",
        raw_description="",
        # caption=get_caption(ct),
    )

    # disabling because I can't get the youtube channel comments insert api endpoint working
    # platform = pb.SocialPlatform.SOCIAL_PLATFORM_YOUTUBE
    # s.posts.append(pb.Post(
    #     platform=platform,
    #     title=f"{s.raw_title}\n\n{get_common_text(ct, platform)}\n\n{get_hashtags(ct, platform)}",
    #     description="N/A",
    #     crosspost=False,
    #     scheduled_unix_timetamp=get_schedule_timestamp(ct),
    # ))

    platform = pb.SocialPlatform.SOCIAL_PLATFORM_INSTAGRAM
    s.posts.append(pb.Post(
        platform=platform,
        title=f"{s.raw_title}\n\n{get_common_text(ct, platform)}\n\n{get_hashtags(ct, platform)}",
        description="N/A",
        crosspost=False,
        scheduled_unix_timetamp=get_schedule_timestamp(ct),
    ))

    # platform = pb.SocialPlatform.SOCIAL_PLATFORM_FACEBOOK
    # s.posts.append(pb.Post(
    #     platform=platform,
    #     title=f"{s.raw_title}\n\n{get_common_text(ct, platform)}\n\n{get_hashtags(ct, platform)}",
    #     description="N/A",
    #     crosspost=False,
    #     scheduled_unix_timetamp=get_schedule_timestamp(ct),
    # ))

    # platform = pb.SocialPlatform.SOCIAL_PLATFORM_TWITTER
    # s.posts.append(pb.Post(
    #     platform=platform,
    #     title=f"{s.raw_title}\n\n{get_common_text(ct, platform)}\n\n{get_hashtags(ct, platform)}",
    #     description="N/A",
    #     crosspost=False,
    #     scheduled_unix_timetamp=get_schedule_timestamp(ct),
    # ))

    return s
