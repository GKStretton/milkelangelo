import typing
import os
import random
from machinepb import machine as pb
import pandas as pd

SOCIAL_TEXT_PATH = "./resources/social_text"

YOUTUBE_SHORT_TITLE_MAX_LENGTH = 90 # actually 100 but that's a bit long

def get_schedule_timestamp(ct: pb.ContentType) -> int:
	# todo: implement
	return 0

def get_random_title_and_description(ct: pb.ContentType) -> typing.Tuple[str, str]:
	return random.choice(load_titles_and_descriptions(ct))

def load_titles_and_descriptions(t: pb.ContentType) -> typing.List[typing.Tuple[str, str]]:
	p = os.path.join(SOCIAL_TEXT_PATH, f"{t.name}.csv")
	with open(p, 'r') as f:
		lines = f.readlines()[1:]
		lines = [line.strip().split("\t") for line in lines]
		return lines

def get_common_text(ct: pb.ContentType, platform: pb.SocialPlatform) -> str:
	# tab separated value file, as csv
	p = os.path.join(SOCIAL_TEXT_PATH, "standard_descriptions", f"{ct.name}.csv")

	# Load the csv file into a pandas DataFrame
	df = pd.read_csv(p, sep='\t')

	# If there's a column matching platform.name, return every cell as a new line
	if platform.name in df.columns:
		lines = df[platform.name].fillna('').tolist()
		return "\n".join(map(str, lines)).strip()  # convert the column to a list and join

	print(f"No column {platform.name} in {p}")
	return ""


def append_title_hashtags(title: str, ct: pb.ContentType, platform: pb.SocialPlatform) -> str:
	max_length = 1000
	if ct == pb.ContentType.CONTENT_TYPE_SHORTFORM and platform == pb.SocialPlatform.SOCIAL_PLATFORM_YOUTUBE:
		max_length = YOUTUBE_SHORT_TITLE_MAX_LENGTH


	tags = get_hashtags_list(ct, platform)

	i = 0
	while i < len(tags) and len(title) + 1 + len(tags[i]) <= max_length:
		title += " " + tags[i]
		i += 1

	return title

def get_hashtags(ct: pb.ContentType, platform: pb.SocialPlatform) -> str:
	return " ".join(get_hashtags_list(ct, platform))

def get_hashtags_list(ct: pb.ContentType, platform: pb.SocialPlatform) -> typing.List[str]:
	# comma separated value file
	p = os.path.join(SOCIAL_TEXT_PATH, "hashtags.csv")

	# Load the csv file into a pandas DataFrame
	df = pd.read_csv(p, sep=',')

	tags = []

	# Append all entries from the "ALL" column
	tags.extend(['#' + tag for tag in df['ALL'].dropna()])

	# If there's a column matching ct.name, add all from that column
	if ct.name in df.columns:
		tags.extend(['#' + random.choice(df[ct.name].dropna())])

	# Add one tag from any columns labeled A-Z
	for letter in list(map(chr, range(65, 91))):  # Generate list of uppercase letters
		if letter in df.columns:
			tags.append('#' + random.choice(df[letter].dropna()))

	# Shuffle the tags
	random.shuffle(tags)

	# Return the tags
	return tags