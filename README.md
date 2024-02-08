# A Study of Light

This is the monorepo for [A Study of Light](https://www.youtube.com/@StudyOfLight).

Backend scope includes:

- [x] Serial-MQTT bridge supporting an event driven architecture
- [x] RTSP streaming for multiple consumers of webcam feeds
- [x] Remote firmware update
- [x] Remote camera crop configuration
- [x] Control interface for system
- [x] Session management system
- [x] OBS integration for live streaming
- [x] Automated video, photo, and state report capture
- [x] Automated post-processing of photos and short- & long-form video content.
- [x] Automated social media posting
- [x] Automated session control

## Instructions

### .env settings

These are the environment variables required

TOP_CAM=/dev/video2 (default)
FRONT_CAM=/dev/video0 (default)
LIGHT_STORES_DIR=/mnt/md0/light-stores/ (default not set, required in .env)

### kv/ settings:

There's a key value store in e.g. `/mnt/md0/light-stores/kv/` which some configuration.

- ENABLE_CONTENT_SCHEDULER_LOOP: controls whether the content scheduler loop will run to automatically upload at regular intervals
- ENABLE_SCHEDULER: controls whether the scheduler will run, which automatically runs sessions each week

Snapshot of kv/ 08/02/2024:

EMAIL_RECIPIENT_MAINTENANCE           MAILJET_API_SECRET          TWITCH_REFRESH_TOKEN  vial-profiles
EMAIL_RECIPIENT_ROUTINE_OPERATIONS    TWITCH_CLIENT_ID            crop_dslr             youtube-credentials.json
EMAIL_RECIPIENT_SOCIAL_NOTIFICATIONS  TWITCH_CLIENT_SECRET        crop_front-cam        youtube_client_secret.json
ENABLE_CONTENT_SCHEDULER_LOOP         TWITCH_EXTENSION_CLIENT_ID  crop_top-cam			OBS_LANDSCAPE_URL
MAILJET_API_KEY                       TWITCH_EXTENSION_SECRET     system-vial-profiles

## Sub-Systems

- [interface](interface/) - A React.JS interface for local control of embedded and backend systems.
- [services](services/)
	- [pygateway](services/pygateway/) - Serial-MQTT bridge for the firmware interface
	- [rtsp](services/rtsp/) - [MediaMTX](https://github.com/aler9/mediamtx) instance for video stream multiplexing
	- [goo](services/goo/) - Bulk of the runtime backend logic is here
	- [dslrcapture](services/dslrcapture/) - Separate service for capturing DSLR images in a loop while session is live.
- [tools](tools/)
	- [content_generation_dslr_edit.py](tools/content_generation_dslr_edit.py) - for post-processing of the dslr images
	- [content_generation_video.py](tools/content_generation_video.py) - for generating video content
	- [content_generation_timelapse.py](tools/content_generation_timelapse.py) - for generating timelapses from all the top-down dslr images

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