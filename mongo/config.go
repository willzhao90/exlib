package mongo

import (
	"github.com/spf13/viper"
)

// Config MongoDb config
type Config struct {
	URI    string `json:"uri"`
	DbName string `json:"dbName"`
}

// GetConfig get MongoDb config
func GetConfig(v *viper.Viper) (*Config, error) {
	var c Config
	err := v.UnmarshalKey("mongo", &c)
	return &c, err
}
