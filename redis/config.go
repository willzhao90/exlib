package redis

import (
	"github.com/spf13/viper"
)

// Config redigo config
type Config struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
}

// GetConfig gets redigo config
func GetConfig(v *viper.Viper) (c Config, err error) {
	err = v.UnmarshalKey("redis", &c)
	return
}
