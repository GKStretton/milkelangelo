# Dark 

Dark is the supporting system monorepo for [A Study of Light](https://www.youtube.com/@StudyOfLight). Also see [Light](https://github.com/GKStretton/Light). This is in a very early alpha. Scope includes:

- [x] Serial-MQTT bridge supporting an event driven architecture
- [x] RTSP streaming for multiple consumers of webcam feeds
- [x] Remote firmware update
- [x] Remote camera crop configuration
- [x] Control interface for system
- [x] Session management system
- [x] OBS integration for live streaming
- [x] Automated video, photo, and state report capture
- [ ] Automated post-processing of short- & long-form video, and photos.
- [ ] Automated social media posting
- [ ] Link to cloud system for remote control

## Instructions


Fill out the .env file in root of repo with your configuration


## Storage

Currently all data is stored to disk at the `-basePath`.

`basePath`:
- session_content/
	- 1/
	- 2/
- session_metadata/
	- 1.yml
	- 2.yml
- session_events/
	- 1.yml
	- 2.yml

Each `session_content` folder has the following format (output content in bold):

- dslr/
	- raw/
	- post/
		- **selected/**
- video/
	- raw/
		- front/
		- top/
		- stills/
	- post/
		- **shortform/**
		- **longform/**
		- stills/
- **stills/**