package config

import (
	"log"

	"github.com/spf13/viper"
)

//InitConfig initialize configuration
func InitConfig() {
	viper.SetDefault("logLevel", "DEBUG")
	viper.SetDefault("workers", 10)

	viper.SetConfigName("order-conf")
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