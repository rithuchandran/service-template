package main

import "github.com/spf13/viper"

func init() {
	viper.AutomaticEnv()
	viper.SetDefault("API_KEY", "**add api key here**")
	viper.SetDefault("SECRET_KEY", "**add secret key here**")
}
