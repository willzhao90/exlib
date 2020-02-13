package payment

import (
	"github.com/spf13/viper"
)

type PoliConfig struct {
	ApiUrl              string `json:"url"`
	MerchangeCode       string `json:"merchant_code"`
	MerchantAuthCode    string `json:"merchant_auth_code"`
	SuccessUrl          string `json:"success_url"`
	FailureUrl          string `json:"failure_url"`
	CancelledUrl        string `json:"cancelled_url"`
	MerchantHomepageURL string `json:"merchant_homepage_url"`
	NotificationUrl     string `json:notification_url`
}

// GetConfig gets poli config
func GetConfig(v *viper.Viper) (PoliConfig, error) {
	var c PoliConfig
	err := v.UnmarshalKey("poli", &c)
	return c, err
}
