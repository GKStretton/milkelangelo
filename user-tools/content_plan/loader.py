import typing
import os
import random
from machinepb import machine as pb

SOCIAL_TEXT_PATH = "./resources/social_text"

YOUTUBE_SHORT_TITLE_MAX_LENGTH = 100

def get_schedule_timestamp(ct: pb.ContentType) -> int:
	# todo: implement
	return 0

def get_random_title_and_description(ct: pb.ContentType) -> typing.Tuple[str, str]:
	return random.choice(load_titles_and_descriptions(ct))

def load_titles_and_descriptions(t: pb.ContentType) -> typing.List[typing.Tuple[str, str]]:
	p = os.path.join(SOCIAL_TEXT_PATH, f"{t.name}.tsv")
	with open(p, 'r') as f:
		lines = f.readlines()[1:]
		lines = [line.strip().split("\t") for line in lines]
		return lines

def get_common(name: str) -> str:
	p = os.path.join(SOCIAL_TEXT_PATH, "common", f"{name}.txt")
	with open(p, 'r') as f:
		return f.read().strip()


def append_title_hashtags(title: str, ct: pb.ContentType, platform: pb.SocialPlatform) -> str:
	tags = get_hashtags_list(ct, platform)
	# todo: add hashtags until title is full
	return title + " " + " ".join(tags)

def get_hashtags(ct: pb.ContentType, platform: pb.SocialPlatform) -> str:
	return " ".join(get_hashtags_list(ct, platform))

def get_hashtags_list(ct: pb.ContentType, platform: pb.SocialPlatform) -> typing.List[str]:
	# todo: load hashtags from file and do selection
	return ["#implementme"]
