import os
import machinepb.machine as pb
from moviepy.editor import VideoClip, AudioClip, AudioFileClip
from moviepy.audio.fx import audio_loop

def load_song(base_dir: str, music_file: str) -> AudioFileClip:
	path = os.path.join(base_dir, "music", music_file)
	return AudioFileClip(path)

def add_music(base_dir: str, content_clip: VideoClip, music_file: str) -> VideoClip:
	audio: AudioClip = audio_loop(load_song(base_dir, music_file), duration=content_clip.duration)
	
	return content_clip.with_audio(audio)
