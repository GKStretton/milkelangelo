
## Content Generation System

### How it works

This diagram may be of use. Also please see the python code in tools/auto_... and tools/videoediting/...

![Diagram](./architecture.dio.svg)

The `for` loop in `__main__` of [auto_video_post.py](tools/auto_video_post.py)
is the most complex part. It handles iteration of state reports in order to get
all the SectionProperties. There is a `delay` concept which delays the property
change, and a `min_duration` concept which forces those properties to persist
for a certain amount of time. Useful for forcing 1x speed for a dispense, because
the dispense state is transient but it should slow down for > 1 second.

Note that SectionProperties will only be added to the ContentDescriptor if they 
are different to the previous props. This reduces the number of subclips, speeding
up generation

### Suggested contributions

This would probably be the easiest way to contribute at this time. There are a 
lot of potential for improvements and extra features. Here are some ideas:

- add visualisations of current state to the videos based on the state reports
	- target ik positions
	- robot position
	- live pipette level indicator graphic;
	- show what colour is being collected etc.
- add textual descriptions of what is happening / commentary
- add a new content type, highlighting some particular aspect
