package kms

import (
	"github.com/spf13/viper"
)

type Config struct {
	Region string
	KmsKey string
}

func GetConfig(v *viper.Viper) (*Config, error) {
	var c Config
	err := v.UnmarshalKey("kms", &c)
	return &c, err
}
