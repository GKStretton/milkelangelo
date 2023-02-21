package config

import "flag"

var (
	RtspHost         = flag.String("rtspHost", "depth", "host for rtsp server")
	TopCamRtspPath   = flag.String("topCamRtspPath", "top-cam", "text after / for top cam stream")
	FrontCamRtspPath = flag.String("frontCamRtspPath", "front-cam", "text after / for front cam stream")
)