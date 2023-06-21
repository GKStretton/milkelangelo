import numpy as np
from moviepy.editor import TextClip, CompositeVideoClip, ImageClip
from moviepy.video.fx.resize import resize


def pulse(t):
    """
    Function to create a pulse effect.

    1. Sine wave
    2. Take absolute value
    3. 1-curve to make it upside down
    4. trim the bottom off with a max function

    This accents the top of the curve and gives a flat period, creating a beat.
    """
    mag = 0.1
    hz = 2
    scale = 1 + mag * max(1-abs(np.sin(hz * np.pi * t)), 0.2)
    return scale


def slow_grow(t):
    return 1+t/30


font = "./resources/fonts/monogram.ttf"
img = "/mnt/md0/light-stores/session_content/60/dslr/post/selected.jpg"

# Define the text
splash_text = "Milk it for all it's worth!"
splash_clip = TextClip(splash_text, fontsize=100, color='yellow', font=font).set_duration(10)

# Apply the pulse effect
pulsing_clip = splash_clip.resize(pulse).rotate(20, resample='bicubic')

# image
i = ImageClip(img).set_duration(splash_clip.duration)
i = i.fx(resize, (900, 900)).fx(resize, slow_grow)

# Create a composite video clip
video = CompositeVideoClip([i.set_position('center'), pulsing_clip.set_position('center')], size=(1920, 1080))

# Write the result to a file
video.write_videofile("splash.mp4", fps=60)
# video.preview()
