from moviepy.video.VideoClip import TextClip
import typing

def build_title(pos: typing.Tuple[int, int], duration: float) -> TextClip:
	text_clip = TextClip("A Study of Light", size=(900, 140), fontsize=115, color='white', font='DejaVu-Serif-Condensed-Italic')
	text_clip = text_clip.set_position(pos)
	text_clip = text_clip.set_duration(duration)

	return text_clip

def build_session_number(metadata: typing.Dict[str, typing.Any], pos: typing.Tuple[int, int], duration: float) -> TextClip:
	session_number_text = f"#{metadata['production_id']}" if metadata['production'] else f"dev#{metadata['id']}"
	text_clip = TextClip(session_number_text, size=(350, 140), fontsize=115, align='west', color='white', font='DejaVu-Serif-Condensed-Italic')
	text_clip = text_clip.set_position(pos)
	text_clip = text_clip.set_duration(duration)

	return text_clip

def build_speed(speed: float, pos: typing.Tuple[int, int], duration: float) -> TextClip:
	speed_text = ""
	if speed == 1.0:
		speed_text = "1x"
	elif speed < 3.0:
		speed_text = f"{speed:.1f}x"
	else:
		speed_text = f"{speed:.0f}x"

	text_clip = TextClip(speed_text, size=(200, 100), fontsize=70, align='west', color='white', font='DejaVu-Serif-Condensed-Italic')
	text_clip = text_clip.set_position(pos)
	text_clip = text_clip.set_duration(duration)

	return text_clip