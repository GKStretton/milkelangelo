package config

import "flag"

var (
	RtspHost         = flag.String("rtspHost", "depth", "host for rtsp server")
	TopCamRtspPath   = flag.String("topCamRtspPath", "top-cam", "text after / for top cam stream")
	FrontCamRtspPath = flag.String("frontCamRtspPath", "front-cam", "text after / for front cam stream")
	SessionBasePath  = flag.String("sessionBasePath", "/mnt/md0/light-stores/sessions", "base path for sessions")
	RawFootagePath   = flag.String("rawFootagePath", "video/raw", "path within session of raw video")
)
