package main

import (
	"log"

	"github.com/Dataman-Cloud/seckilling/queue/src/handler"

	"github.com/spf13/viper"
	fsnotify "gopkg.in/fsnotify.v1"
)

func initConfig() {
	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", ":5090")
	viper.SetDefault("logLevel", "DEBUG")
	viper.SetDefault("watch", false)
	viper.SetDefault("cache.poolSize", 100)

	viper.SetConfigName("queue-conf")
	viper.SetConfigType("yaml")
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
		handler.ResetSeckillTime()
	})

	log.Printf("loading config %s \n", viper.ConfigFileUsed())
}
