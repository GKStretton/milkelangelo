import typing
import yaml
import os
import random
from machinepb import machine as pb
import pandas as pd
from datetime import datetime, timedelta
import pytz
from betterproto import Casing

SOCIAL_TEXT_PATH = "../resources/social_text"

YOUTUBE_SHORT_TITLE_MAX_LENGTH = 90  # actually 100 but that's a bit long


def load_content_statuses(base_dir: str, session_number: int) -> typing.Optional[pb.ContentTypeStatuses]:
    path = os.path.join(base_dir, "session_content", session_number, "content_plan.yml")
    try:
        with open(path, 'r') as f:
            yml = yaml.load(f, Loader=yaml.FullLoader)
            meta: pb.ContentTypeStatuses = pb.ContentTypeStatuses().from_dict(value=yml)
            return meta
    except FileNotFoundError:
        return None


def write_content_statuses(content_statuses: pb.ContentTypeStatuses, base_dir: str, session_number: int):
    output_path = os.path.join(base_dir, "session_content", session_number, "content_plan.yml")
    d = content_statuses.to_dict(casing=Casing.SNAKE, include_default_values=True)
    with open(output_path, 'w') as f:
        yaml.dump(d, f)


def get_schedule_timestamp(ct: pb.ContentType) -> int:
    schedule_hour = 18
    schedule_minute = 30

    now = datetime.now(pytz.utc)

    # schedule for the next day at schedule_hour
    next_schedule = (now + timedelta(days=1)).replace(hour=schedule_hour,
                                                      minute=schedule_minute, second=0, microsecond=0)

    if ct == pb.ContentType.CONTENT_TYPE_STILL:
        return 0
    elif ct == pb.ContentType.CONTENT_TYPE_LONGFORM:
        # 1 day later
        # return int((next_schedule).timestamp())
        # immediate because why not?
        return 0
    elif ct == pb.ContentType.CONTENT_TYPE_SHORTFORM:
        # 2 days later
        # return int((next_schedule + timedelta(days=1)).timestamp())
        # immediate because why not?
        return 0
    elif ct == pb.ContentType.CONTENT_TYPE_CLEANING:
        # Immediate because unlisted
        return 0
    elif ct == pb.ContentType.CONTENT_TYPE_DSLR:
        # Immediate because unlisted
        return 0
    else:
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


def get_caption(ct: pb.ContentType) -> str:
    # tab separated value file
    p = os.path.join(SOCIAL_TEXT_PATH, "captions.csv")

    # Load the csv file into a pandas DataFrame
    df = pd.read_csv(p, sep='\t')

    captions = list(df['ALL'].dropna())
    if ct == pb.ContentType.CONTENT_TYPE_LONGFORM or ct == pb.ContentType.CONTENT_TYPE_SHORTFORM:
        captions.extend(list(df['CONTENT_TYPE_SHORTFORM'].dropna()))

    if ct == pb.ContentType.CONTENT_TYPE_LONGFORM or ct == pb.ContentType.CONTENT_TYPE_CLEANING:
        captions.extend(list(df['CONTENT_TYPE_CLEANING'].dropna()))

    return random.choice(captions)


def append_title_hashtags(title: str, ct: pb.ContentType, platform: pb.SocialPlatform) -> str:
    max_length = 1000
    if (ct == pb.ContentType.CONTENT_TYPE_SHORTFORM or ct == pb.ContentType.CONTENT_TYPE_CLEANING or ct == pb.ContentType.CONTENT_TYPE_DSLR) and platform == pb.SocialPlatform.SOCIAL_PLATFORM_YOUTUBE:
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
    # tab separated value file
    p = os.path.join(SOCIAL_TEXT_PATH, "hashtags.csv")

    # Load the csv file into a pandas DataFrame
    df = pd.read_csv(p, sep='\t')

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


def get_splashtext() -> str:
    # comma separated value file
    p = os.path.join(SOCIAL_TEXT_PATH, "splashtext.csv")

    # Load the csv file into a pandas DataFrame
    df = pd.read_csv(p, sep='\t')

    return random.choice(df["TEXT"].dropna())


def generate_splashtext_hue():
    """
    generate a hue for the splashtext, excluding yellow and dark blue
    """
    choices = []
    for i in range(0, 360):
        if i >= 50 and i <= 80:
            continue
        if i >= 225 and i <= 260:
            continue
        choices.append(i)
    return random.choice(choices)
