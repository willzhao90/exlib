package config

import (
	"fmt"

	"github.com/spf13/viper"
)

/*
* Load config for service based on env
* Local: load a local yaml
* Dev or Prod: env variables
 */
func LoadConfig(svc string) (v *viper.Viper, err error) {
	v = viper.New()
	v.SetConfigFile(fmt.Sprintf("config/%s.yaml", svc))
	err = v.ReadInConfig()
	// v = viper.New()
	// pflag.StringP("env", "e", "local", "env where this service is running")
	// pflag.Parse()
	// viper.BindPFlags(pflag.CommandLine)

	// env := v.GetString("env")
	// todo, a better way to load local config, can be same as
	// the other two.
	// v.SetConfigFile(fmt.Sprintf("config/%s.yaml", svc))
	// err = v.ReadInConfig()
	// log.Println(env)
	// if env == "" || env == "local" {
	// 	// todo, a better way to load local config, can be same as
	// 	// the other two.
	// 	v.SetConfigFile("config/local.yaml")
	// 	err = v.ReadInConfig()
	// } else {
	// 	v.SetConfigFile(fmt.Sprintf("%s.yaml", svc))
	// 	err = v.ReadInConfig()
	// }
	return
}

// GetLocalConfig gets viper
func GetLocalConfig(path string) (v *viper.Viper, err error) {
	v = viper.New()
	v.SetConfigFile(path)
	err = v.ReadInConfig()
	return
}
