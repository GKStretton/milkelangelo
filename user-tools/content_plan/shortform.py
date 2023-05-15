from machinepb import machine as pb
from content_plan.loader import *

def buildShortform(n: int) -> pb.ContentTypeStatus:
	ct = pb.ContentType.CONTENT_TYPE_SHORTFORM
	raw_title, raw_description = get_random_title_and_description(ct)

	s = pb.ContentTypeStatus(
		raw_title=raw_title,
		raw_description=raw_description,
		caption="define captions!",
	)

	platform = pb.SocialPlatform.SOCIAL_PLATFORM_YOUTUBE
	s.posts.append(pb.Post(
		platform=platform,
		title=append_title_hashtags(f"{s.raw_title} - {n}", ct, platform),
		description=f"{get_hashtags(ct, platform)}\n\n{s.raw_description}\n\n{get_common('description_shortform_youtube')}",
		crosspost=False,
		scheduled_unix_timetamp=get_schedule_timestamp(ct),
	))

	return s
