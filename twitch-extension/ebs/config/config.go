package config

import (
	"github.com/spf13/viper"
)

var defaults = map[string]any{
	"SHARED_SECRET_GOO":            "local_secret",
	"SHARED_SECRET_TWITCH":         "",
	"BROADCAST_STATE_TO_TWITCH":    false,
	"ENABLE_SERVER_AUTHENTICATION": false,
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

func SharedSecretGoo() string {
	return viper.GetString("SHARED_SECRET_GOO")
}

func SharedSecretTwitch() string {
	return viper.GetString("SHARED_SECRET_TWITCH")
}

func BroadcastStateToTwitch() bool {
	return viper.GetBool("BROADCAST_STATE_TO_TWITCH")
}

func EnableServerAuthentication() bool {
	return viper.GetBool("ENABLE_SERVER_AUTHENTICATION")
}
