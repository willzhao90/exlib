package sms

import (
	"github.com/spf13/viper"
)

// Config twilio config
type Config struct {
	AccountSid string `json:"accountSid"`
	AuthToken  string `json:"authToken"`
	URL        string `json:"url"`
	From       string `json:"from"`
}

// GetTwilioConfig get twilio config
func GetTwilioConfig(v *viper.Viper) (*Config, error) {
	var c Config
	err := v.UnmarshalKey("twilio", &c)
	return &c, err
}
