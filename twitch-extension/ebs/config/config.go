package config

import (
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Test string `env:"TEST"`
}

var defaults = &Config{
	Test: "Banana",
}

func setDefaults(v *viper.Viper, defaults interface{}) {
	val := reflect.ValueOf(defaults).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		key := strings.ToLower(field.Name)
		v.SetDefault(key, value.Interface())
	}
}

func GetConfig() *Config {
	v := viper.NewWithOptions(viper.ExperimentalBindStruct())
	v.AutomaticEnv()
	setDefaults(v, defaults)

	c := &Config{}
	v.Unmarshal(&c)

	return c
}
