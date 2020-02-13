package pusher

import (
	"github.com/spf13/viper"
)

// Config required by pusher
type Config struct {
	AppID   string `json:"appId"`
	Key     string `json:"key"`
	Secret  string `json:"secret"`
	Cluster string `json:"cluster"`
	Secure  bool   `json:"secure"`
}

// GetConfig fetches the kafka config from viper
func GetConfig(v *viper.Viper) (Config, error) {
	var c Config
	err := v.UnmarshalKey("pusher", &c)
	return c, err
}
