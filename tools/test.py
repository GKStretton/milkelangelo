from videoediting.intro import build_intro
from videoediting.outro import build_outro
import machinepb.machine as pb
from videoediting.loaders import get_session_metadata, get_content_plan
from moviepy.video.fx import resize


def intros():
    """Manually generate and preview an intro"""
    session_number = 59
    content_type = pb.ContentType.CONTENT_TYPE_LONGFORM

    bd = "/mnt/md0/light-stores"
    metadata = get_session_metadata(bd, session_number)
    content_plan = get_content_plan(bd, session_number)

    intro = build_intro(bd, session_number, metadata, content_type, content_plan, 3.33)
    # intro.write_videofile("intro.mp4", fps=60)
    intro.fx(resize, 0.5).preview()


def outros():
    """Manually generate and preview an intro"""
    session_number = 60
    content_type = pb.ContentType.CONTENT_TYPE_LONGFORM

    bd = "/mnt/md0/light-stores"
    metadata = get_session_metadata(bd, session_number)

    outro = build_outro(bd, session_number, metadata, content_type, 13.33)
    # outro.write_videofile("outro.mp4", fps=60)
    outro.fx(resize, 0.5).preview()


if __name__ == "__main__":
    intros()
