package auth0

import (
	"github.com/spf13/viper"
)

// Config auth0 config
type Config struct {
	URL          string `json:"url"`
	Domain       string `json:"domain"`
	Audience     string `json:"auth0_audience"`
	Connection   string `json:"connection"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	ApiToken     string `json:"api_token"`
}

// GetConfig gets auth0 config
func GetConfig(v *viper.Viper) (Config, error) {
	var c Config
	err := v.UnmarshalKey("auth", &c)
	return c, err
}
