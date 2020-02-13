package withdraw

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

type WithdrawConfig struct {
	WithdrawFee map[string]map[string]string
}

// GetConfig ..
func GetConfig(v *viper.Viper) (*WithdrawConfig, error) {
	conf := new(WithdrawConfig)
	conf.WithdrawFee = map[string]map[string]string{}
	for _, k := range v.AllKeys() {
		arr := strings.Split(k, ".")
		if arr == nil || len(arr) != 2 {
			return nil, errors.New("Invalid entry in withdraw config file.")
		}
		key := arr[0]
		conf.WithdrawFee[key] = map[string]string{}
		mp := v.Get(key).(map[string]interface{})
		var ok bool
		for x, y := range mp {
			conf.WithdrawFee[key][x], ok = y.(string)
			if !ok {
				return nil, errors.New("Invalid entry in withdraw config file.")
			}
		}
	}
	return conf, nil
}
