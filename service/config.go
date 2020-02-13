package service

import (
	"github.com/spf13/viper"
)

// Config contains addresses of services
type Config struct {
	Trading    string `json:"trading"`    // service.trading
	Member     string `json:"member"`     // service.member
	Currency   string `json:"currency"`   // service.currency
	Otc        string `json:"otc"`        // service.otc
	Wallet     string `json:"wallet"`     // service.wallet
	FailLoader string `json:"failLoader"` // service.fail_loader
	ID3Global  string `json:"id3Global"`  // service.id3global
	Admin      string `json:"admin"`      // service.admin
	Bridge     string `json:"bridge"`     // service.bridge
	AutoQ      string `json:"autoq"`      // service.autoq
} // todo fixme

// GetConfig gets grpc server addresses
func GetConfig(v *viper.Viper) (Config, error) {
	var c Config
	err := v.UnmarshalKey("service", &c)
	return c, err
}
