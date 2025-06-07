# Milkelangelo

_Previously "A Study of Light"_

This is the monorepo for [Milkelangelo](https://www.youtube.com/@StudyOfLight).

- 3d files
- Firmware
- serial-mqtt bridge with remote firmware updates
- protos and constants
- Backend Systems
	- [x] RTSP streaming for multiple consumers of webcam feeds
	- [x] Web control interface for system
	- [x] "Session" management
	- [x] OBS integration for live streaming
	- [x] Automated video, photo, and state report capture
	- [x] Automated post-processing of photos and short- & long-form video content.
	- [x] Automated social media posting
	- [x] Automated session control

## Usage

### Installation

- Install Ubuntu on a machine
- Disable automounting in gnome:

```bash
gsettings set org.gnome.desktop.media-handling automount false
gsettings set org.gnome.desktop.media-handling automount-open false
```

- [Install docker](https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository). ** NOT "docker desktop" ** as it doesn't work nicely with /dev/ devices (serial and webcam)
- Clone repo
- Make .env file (see below for settings)
- `sudo docker compose up -d`
- Install OBS
  - install gstreamer plugin
  - profile and scene collections are in `resources/obs`
  - enable websocket server

### scripts

- `./scripts/dev/mount-md0` - mount milkelangelo's filesystem locally, /mnt/md0/light-stores
- `./scripts/dev/extras-update` - render protos and xlangconsts, copying to everywhere necessary.
- `./scripts/dev/update-machine` - update milkelangelo minipc with latest version of the backend systems
- `./scripts/dev/cloud-certs` - update webrtc.gregstretton.org certs for when the webcam stream fails on the cloud system
- `./scripts/dev/cloud-update` - update asol.gregstretton.org interface with latest build from github action

### .env settings

These are the environment variables used, and they default to the following in the docker-compose:

TOP_CAM=/dev/video2
FRONT_CAM=/dev/video0
LIGHT_STORES_DIR=/mnt/md0/light-stores/

### kv/ settings:

There's a key-value store in e.g. `/mnt/md0/light-stores/kv/` with most configuration.

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

## License

This work is released under the CC0 1.0 Universal Public Domain Dedication. You can find the full text of the license here: https://creativecommons.org/publicdomain/zero/1.0/

### Polite Request for Attribution

While it's not legally required, we kindly ask that you give credit if you use or modify this work. Attribution helps support the project and encourages future learning and contributions. You can provide credit by linking to this repository or mentioning the original author's name. Thank you for your cooperation!