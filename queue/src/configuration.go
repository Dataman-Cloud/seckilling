package main

import (
	"log"

	"github.com/spf13/viper"
	"gopkg.in/fsnotify.v1"
)

func initConfig() {
	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", ":5090")
	viper.SetDefault("logLevel", "DEBUG")
	viper.SetDefault("watch", false)

	viper.SetConfigName("queue-conf")
	viper.SetConfigType("json")
	viper.AddConfigPath(".") // optionally look for config in the working directory
	viper.AddConfigPath("$HOME/.seckilling/")
	viper.AddConfigPath("/etc/seckilling/")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Panicf("Fatal error config file: %s \n", err)
	}

	if viper.GetBool("watch") {
		viper.WatchConfig()
	}
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
	})

	log.Printf("loading config %s \n", viper.ConfigFileUsed())
}
