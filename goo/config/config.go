package config

// todo
// currently config is a mess combination of envvar, flags, and keyvalue store

// requirements
// 1. typed configuration via this package, so everything is got via here. flat structure probably
// 2. support for default values per environment - local, and deployed
// 3. support for overriding values.
// 4. validation and documentation is a bonus

// maybe a struct like this that is just public to avoid all the getter stuff.
type Config struct {
	EbsHost string `json:"ebs_host"`
}

func GetConfig() *Config {
	return nil
}
