import argparse
from machinepb import machine as pb
from betterproto import Casing
import yaml
import os

def buildCleaning() -> pb.ContentTypeStatus:
	return None

def buildDslr() -> pb.ContentTypeStatus:
	return None

def buildShortform() -> pb.ContentTypeStatus:
	return None

def buildLongform() -> pb.ContentTypeStatus:
	longform = pb.ContentTypeStatus(
		raw_title="this is a title",
		raw_description="this is a desc",
		caption="watch this!",
	)

	longform_youtube = pb.Post(
		platform=pb.SocialPlatform.SOCIAL_PLATFORM_YOUTUBE,
		title=longform.raw_title+" YT",
		description=longform.raw_description+" YT",
		crosspost=False,
		url="blah"
	)

	longform.posts.append(longform_youtube)

	return longform


if __name__ == "__main__":
	parser = argparse.ArgumentParser()
	parser.add_argument("-d", "--base-dir", action="store", help="base directory containing session_content and session_metadata", required=True)
	parser.add_argument("-n", "--session-number", action="store", help="session number e.g. 5", required=True)
	args = parser.parse_args()
	output_path = os.path.join(args.base_dir, "session_content", args.session_number, "content_plan.yml")

	print(f"Launching publish_helpers for session {args.session_number} in '{args.base_dir}'\n")

	content_statuses = pb.ContentTypeStatuses()

	content_statuses.content_statuses[pb.ContentType.CONTENT_TYPE_CLEANING.name] = buildCleaning()
	content_statuses.content_statuses[pb.ContentType.CONTENT_TYPE_DSLR.name] = buildDslr()
	content_statuses.content_statuses[pb.ContentType.CONTENT_TYPE_LONGFORM.name] = buildLongform()
	content_statuses.content_statuses[pb.ContentType.CONTENT_TYPE_SHORTFORM.name] = buildShortform()

	d = content_statuses.to_dict(casing=Casing.SNAKE, include_default_values=True)
	with open(output_path, 'w') as f:
		yaml.dump(d, f)