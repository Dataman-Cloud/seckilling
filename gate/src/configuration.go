package main

import (
	"log"

	"github.com/spf13/viper"
	fsnotify "gopkg.in/fsnotify.v1"
)

func initConfig() {
	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", ":5090")
	viper.SetDefault("logLevel", "DEBUG")
	viper.SetDefault("watch", false)
	viper.SetDefault("cache.poolSize", 100)

	viper.BindEnv("init_model", "INIT_MODEL")
	model := viper.GetInt("init_model")
	switch model {
	case 0:
		initByConfigFile()
	case 1:
		initByEnv()
	default:
		initByConfigFile()
	}
}

func initByEnv() {
	log.Println("init params by env..")
	viper.BindEnv("host", "HOST")
	viper.BindEnv("port", "PORT")
	viper.BindEnv("logLevel", "LOG_LEVEL")
	viper.BindEnv("cache.host", "CACHE_HOST")
	viper.BindEnv("cache.port", "CACHE_PORT")
	viper.BindEnv("cache.poolSize", "CACHE_POOLSIZE")
}

func initByConfigFile() {
	log.Println("init params by configfile..")
	viper.SetConfigName("gate-conf")
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
	})

	log.Printf("loading config %s \n", viper.ConfigFileUsed())
}
