import argparse
from machinepb import machine as pb
from betterproto import Casing
import yaml
import os

from content_plan.cleaning import *
from content_plan.dslr import *
from content_plan.longform import *
from content_plan.shortform import *


if __name__ == "__main__":
	parser = argparse.ArgumentParser()
	parser.add_argument("-d", "--base-dir", action="store", help="base directory containing session_content and session_metadata", required=True)
	parser.add_argument("-n", "--session-number", action="store", help="session number e.g. 5", required=True)
	args = parser.parse_args()

	output_path = os.path.join(args.base_dir, "session_content", args.session_number, "content_plan.yml")
	n = args.session_number

	print(f"Launching publish_helpers for session {args.session_number} in '{args.base_dir}'\n")

	content_statuses = pb.ContentTypeStatuses()

	content_statuses.content_statuses[pb.ContentType.CONTENT_TYPE_LONGFORM.name] = buildLongform(n)
	content_statuses.content_statuses[pb.ContentType.CONTENT_TYPE_SHORTFORM.name] = buildShortform(n)
	# content_statuses.content_statuses[pb.ContentType.CONTENT_TYPE_CLEANING.name] = buildCleaning(n)
	# content_statuses.content_statuses[pb.ContentType.CONTENT_TYPE_DSLR.name] = buildDslr(n)
	# content_statuses.content_statuses[pb.ContentType.CONTENT_TYPE_STILL.name] = buildStill(n)

	d = content_statuses.to_dict(casing=Casing.SNAKE, include_default_values=True)
	with open(output_path, 'w') as f:
		yaml.dump(d, f)