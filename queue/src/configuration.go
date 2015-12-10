package main

import (
	"log"

	"github.com/spf13/viper"
)

func initConfig() {
	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", ":5090")
	viper.SetDefault("logLevel", "DEBUG")

	viper.SetConfigName("queue-conf")
	viper.SetConfigType("json")
	viper.AddConfigPath(".") // optionally look for config in the working directory
	viper.AddConfigPath("$HOME/.seckilling/")
	viper.AddConfigPath("/etc/seckilling/")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Panicf("Fatal error config file: %s \n", err)
	}
	log.Printf("loading config %s \n", viper.ConfigFileUsed())
}
