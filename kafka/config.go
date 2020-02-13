package kafka

import (
	"github.com/spf13/viper"
)

// Config required by kafka
type Config struct {
	Brokers []string `json:"brokers"`
	// DispatcherErrFatal: true: consumer dispatcher fail and break loop when error occurred
	DispatcherErrFatal bool `json:"dispatcherErrFatal"`
}

// GetConfig fetches the kafka config from viper
func GetConfig(v *viper.Viper) (*Config, error) {
	var c Config
	err := v.UnmarshalKey("kafka", &c)
	return &c, err
}
