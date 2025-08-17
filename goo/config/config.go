package config

import (
	"github.com/spf13/viper"
)

var defaults = map[string]any{
	"LIGHT_STORES_DIR":  "/mnt/md0/light-stores/",
	"BROKER_HOST":       "milkelangelo",
	"ENABLE_EBS":        true,
	"EBS_HOST":          "localhost",
	"SHARED_SECRET_EBS": "local_secret",
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

func EnableEBS() bool {
	return viper.GetBool("ENABLE_EBS")
}

func EbsHost() string {
	return viper.GetString("EBS_HOST")
}

func SharedSecretEbs() string {
	return viper.GetString("SHARED_SECRET_EBS")
}
