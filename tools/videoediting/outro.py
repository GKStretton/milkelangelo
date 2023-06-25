import machinepb.machine as pb
from moviepy.editor import ColorClip, CompositeVideoClip, TextClip
from videoediting.constants import Format, MAIN_FONT, PIXEL_FONT
from videoediting.intro import build_base
from videoediting.compositor_helpers import (
    get_size_from_format,
    build_dslr_image,
    slow_grow_effect
)


def timezones():
    from datetime import datetime, time
    from dateutil.relativedelta import relativedelta
    import pytz

    # Create a datetime object for 18:30 UTC
    dt_utc = datetime.combine(datetime.today() + relativedelta(weeks=1), time(18, 30)).replace(tzinfo=pytz.UTC)

    # Convert to UK time
    zones = [
        'America/Los_Angeles',
        'America/New_York',
        'America/Sao_Paulo',
        'UTC',
        'Europe/London',
        'Europe/Berlin',
        # 'Asia/Shanghai',
        'Asia/Kolkata',
        # 'Australia/Sydney',
        # 'Africa/Johannesburg',
        # 'Asia/Tokyo',
    ]

    outputs = []
    for z in zones:
        dt = dt_utc.astimezone(pytz.timezone(z))

        # Format the output
        time = dt.strftime('%H:%M')
        zone = dt.strftime('%Z')
        if zone == '-03':
            zone = "BRT"
        output = (time, zone)
        print(output)
        outputs.append(output)

    return outputs


def build_outro(
    base_dir: str,
    session_number: int,
    metadata,
    content_type: pb.ContentType,
    duration: int
):
    fmt = Format.PORTRAIT
    if content_type == pb.ContentType.CONTENT_TYPE_LONGFORM:
        fmt = Format.LANDSCAPE

    title, session_number_clip, _ = build_base(
        metadata,
        duration,
        "unused",
        portrait=fmt == Format.PORTRAIT
    )
    dslr = build_dslr_image(
        base_dir,
        session_number,
        duration,
        fmt,
        ('center'
         if fmt == Format.PORTRAIT else
         lambda t: (870 + (1-slow_grow_effect(t))*500, 'center')),
    )

    clips = [dslr, title, session_number_clip]

    # Date
    clips.append(
        TextClip(
            "Join LIVE every SATURDAY\n\n\n\ntwitch.tv/StudyOfLight",
            color='white',
            font_size=70,
            font=PIXEL_FONT
        )
        .with_position((200, 390))
        .with_duration(duration)
    )

    for i, t in enumerate(timezones()):
        print(i, t)
        clips.append(
            TextClip(
                f"{t[1]}\n{t[0]}",
                color='white',
                font_size=60,
                font=PIXEL_FONT
            )
            .with_position((440 + (i-3)*140, 500))
            .with_duration(duration)
        )

    return CompositeVideoClip(clips, size=get_size_from_format(fmt))
