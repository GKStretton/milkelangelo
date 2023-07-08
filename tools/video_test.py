from moviepy.editor import VideoFileClip
from moviepy.video.fx import resize, accel_decel
import math

vid = VideoFileClip("/mnt/md0/light-stores/session_content/60/video/post/CONTENT_TYPE_SHORTFORM.7.mp4")
vid = vid.fx(resize, 0.5)
vid = vid.time_transform(lambda t: math.sin(t))
vid = vid.with_duration(20)
vid.preview()
