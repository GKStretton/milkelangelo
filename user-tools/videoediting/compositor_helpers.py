from moviepy.video.VideoClip import TextClip
import typing

def build_title(pos: typing.Tuple[int, int], duration: float) -> TextClip:
	text_size=(900, 140)
	text_clip = TextClip("A Study of Light", size=text_size, fontsize=115, color='white', font='DejaVu-Serif-Condensed-Italic')
	text_clip = text_clip.set_position((pos[0] - text_size[0] // 2, pos[1]))
	text_clip = text_clip.set_duration(duration)

	return text_clip

def build_session_number(metadata: typing.Dict[str, typing.Any], pos: typing.Tuple[int, int], duration: float, center: bool = False) -> TextClip:
	text_size=(350, 140)
	session_number_text = f"#{metadata['production_id']}" if metadata['production'] else f"dev#{metadata['id']}"
	text_clip = TextClip(session_number_text, size=text_size, fontsize=115, align='center' if center else 'west', color='white', font='DejaVu-Serif-Condensed-Italic')
	text_clip = text_clip.set_position((pos[0] - text_size[0] // 2, pos[1]))
	text_clip = text_clip.set_duration(duration)

	return text_clip

def build_subtitle(text: str, pos: typing.Tuple[int, int], duration: float) -> TextClip:
	text_size=(1080, 400)
	text_clip = TextClip(text, size=text_size, fontsize=80, align='center', color='white', font='DejaVu-Serif-Condensed-Italic')
	text_clip = text_clip.set_position((pos[0] - text_size[0] // 2, pos[1]))
	text_clip = text_clip.set_duration(duration)

	return text_clip

def closest_rounded_speed(speed: float) -> str:
    valid_speeds = [0.5, 1, 1.5, 2, 2.5, 3, 4, 5, 10, 15, 20, 25, 30, 40, 50, 75, 100, 150, 200, 300, 400, 500, 750, 1000]
	# minimum difference between speed and one of valid_speeds
    closest_speed = min(valid_speeds, key=lambda x: abs(x - speed))

    if closest_speed == 1 or closest_speed == 2:
        return f"{closest_speed:.0f}x"
    elif closest_speed < 3.0:
        return f"{closest_speed:.1f}x"
    else:
        return f"{closest_speed:.0f}x"

def build_speed(speed: float, pos: typing.Tuple[int, int], duration: float) -> TextClip:
	speed_text = closest_rounded_speed(speed)

	text_clip = TextClip(speed_text, size=(200, 100), fontsize=70, align='west', color='white', font='DejaVu-Serif-Condensed-Italic')
	text_clip = text_clip.set_position(pos)
	text_clip = text_clip.set_duration(duration)

	return text_clip