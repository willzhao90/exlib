package alicloud

import (
	"github.com/spf13/viper"
)

type AliConfig struct {
	Bucket    string `json:"bucket"`
	Region    string `json:"region"`
	Endpoint  string `json:"endpoint"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

// GetConfig gets auth0 config
func GetConfig(v *viper.Viper) (AliConfig, error) {
	var c AliConfig
	err := v.UnmarshalKey("ali", &c)
	return c, err
}
