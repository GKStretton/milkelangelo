import argparse
from machinepb import machine as pb
import yaml
import os
import videoediting.loaders as loaders

from content_plan.cleaning import *
from content_plan.dslr import *
from content_plan.longform import *
from content_plan.shortform import *
from content_plan.still import *

from content_plan.loader import *

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("-d", "--base-dir", action="store",
                        help="base directory containing session_content and session_metadata", default="/mnt/md0/light-stores")
    parser.add_argument("-n", "--session-number", action="store", help="session number e.g. 5", required=True)
    parser.add_argument("-o", "--override", action="store_true",
                        help="if true, override rather than quitting if already exists", default=False)
    args = parser.parse_args()

    metadata = loaders.get_session_metadata(args.base_dir, args.session_number)
    session_number_text = f"#{metadata['production_id']}" if metadata['production'] else f"dev#{metadata['id']}"

    print(f"Launching content plan generator for session {args.session_number} in '{args.base_dir}'\n")

    content_statuses = load_content_statuses(args.base_dir, args.session_number)
    if content_statuses is not None and not args.override:
        print("content statuses already exist for this session. Use -o / --override to overwrite.")
        exit(0)

    content_statuses = pb.ContentTypeStatuses()

    content_statuses.splashtext = get_splashtext()
    content_statuses.splashtext_hue = generate_splashtext_hue()

    content_statuses.content_statuses[pb.ContentType.CONTENT_TYPE_LONGFORM.name] = buildLongform(args.base_dir, session_number_text)
    content_statuses.content_statuses[pb.ContentType.CONTENT_TYPE_SHORTFORM.name] = buildShortform(args.base_dir, session_number_text)
    content_statuses.content_statuses[pb.ContentType.CONTENT_TYPE_CLEANING.name] = buildCleaning(args.base_dir, session_number_text, unlisted=True)
    content_statuses.content_statuses[pb.ContentType.CONTENT_TYPE_DSLR.name] = buildDslr(session_number_text, unlisted=True)
    content_statuses.content_statuses[pb.ContentType.CONTENT_TYPE_STILL.name] = buildStill(session_number_text)

    write_content_statuses(content_statuses, args.base_dir, args.session_number)
