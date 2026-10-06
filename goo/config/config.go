package config

import (
	"github.com/spf13/viper"
)

var defaults = map[string]any{
	"LIGHT_STORES_DIR":  "/mnt/md0/light-stores/",
	"BROKER_HOST":       "milkelangelo",
	"ENABLE_PUBLIC_API": true,
	"PUBLIC_API_ADDR":   "127.0.0.1:8789",
	"PUBLIC_UI_DIR":     "",
	"MEDIAMTX_URL":      "http://milkelangelo:8889",
}

func init() {
	viper.AutomaticEnv()
	setDefaults(defaults)
}

func setDefaults(defaults map[string]any) {
	for key, value := range defaults {
		viper.SetDefault(key, value)
	}
}

func LightStores() string {
	return viper.GetString("LIGHT_STORES_DIR")
}

func BrokerHost() string {
	return viper.GetString("BROKER_HOST")
}

func EnablePublicApi() bool {
	return viper.GetBool("ENABLE_PUBLIC_API")
}

func PublicApiAddr() string {
	return viper.GetString("PUBLIC_API_ADDR")
}

// PublicUiDir is the built remote control page, served by the public api.
func PublicUiDir() string {
	return viper.GetString("PUBLIC_UI_DIR")
}

// MediamtxURL is MediaMTX's WebRTC server, whose signalling for the bowl
// camera is proxied by the public api. Empty disables video.
func MediamtxURL() string {
	return viper.GetString("MEDIAMTX_URL")
}
