# A Study of Light - Backend Systems

This is the supporting system monorepo for [A Study of Light](https://www.youtube.com/@StudyOfLight). Also see [Light](https://github.com/GKStretton/Light). Scope includes:

- [x] Serial-MQTT bridge supporting an event driven architecture
- [x] RTSP streaming for multiple consumers of webcam feeds
- [x] Remote firmware update
- [x] Remote camera crop configuration
- [x] Control interface for system
- [x] Session management system
- [x] OBS integration for live streaming
- [x] Automated video, photo, and state report capture
- [x] Automated post-processing of photos and short- & long-form video content.
- [ ] Automated social media posting

## Instructions

Fill out the .env file in root of repo with your configuration. This wasn't initially
built to be used by others so let me know any problems or missing information!

## Sub-Systems

- [interface](interface/) - A React.JS interface for local control of embedded and backend systems.
- [services](services/)
	- [pygateway](services/pygateway/) - Serial-MQTT bridge for the firmware interface
	- [rtsp](services/rtsp/) - [MediaMTX](https://github.com/aler9/mediamtx) instance for video stream multiplexing
	- [goo](services/goo/) - Bulk of the runtime backend logic is here
	- [dslrcapture](services/dslrcapture/) - Separate service for capturing DSLR images in a loop while session is live.
- [tools](user-tools/)
	- [auto_image_post.py](user-tools/auto_image_post.py) - for post-processing of the dslr images
	- [auto_stills_editing.py](user-tools/auto_stills_editing.py) - for generating "thumbnails"
	- [auto_video_post.py](user-tools/auto_video_post.py) - for generating video content
	- [auto_dslr_timelapse.py](user-tools/auto_dslr_timelapse.py) - for generating timelapses from all the top-down dslr images

## Storage

Currently all data is stored to disk at the `-basePath` provided to _goo_.

`basePath`:
- session_content/
	- 1/
	- 2/
- session_metadata/
	- 1.yml
	- 2.yml
- kv/ (key-value storage for crop config etc.)
	- crop_dslr
	- crop_front-cam
	- crop_top-cam

Each `session_content` folder has the following format (output content in bold):

- state-reports.yml
- dispense-metadata.yml
- dslr/
	- raw/
		- 0001.jpg
		- 0001.jpg.creationtime
		- 0001.jpg.yml (crop info)
	- post/
		- **selected/**
- video/
	- raw/
		- front-cam/
		- top-cam/
	- post/
		- **CONTENT_TYPE.i.mp4**
- stills/
	- **INTRO/OUTRO-LANDSCAPE/PORTRAIT.jpg**

## Content Generation System

### How it works

This diagram may be of use. Also please see the python code in user-tools/auto_... and user-tools/videoediting/...

![Diagram](./architecture.dio.svg)

The `for` loop in `__main__` of [auto_video_post.py](user-tools/auto_video_post.py)
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

## Contributing

It would be fantastic if you'd like to contribute to this project, there's so many
possible additions to this project that it'd be impossible for 1 person to ever
max out this potential.

Please submit an issue with ideas and we can get in sync.
I can provide a session's data and raw footage for development purposes if helpful.

Then once scope is agreed, make your contribution and submit a pull request!

## License

This work is released under the CC0 1.0 Universal Public Domain Dedication. You can find the full text of the license here: https://creativecommons.org/publicdomain/zero/1.0/

### Polite Request for Attribution

While it's not legally required, we kindly ask that you give credit if you use or modify this work. Attribution helps support the project and encourages future learning and contributions. You can provide credit by linking to this repository or mentioning the original author's name. Thank you for your cooperation!