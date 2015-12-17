package main

import (
	"fmt"
	"log"

	"github.com/Dataman-Cloud/seckilling/order/src/config"
	"github.com/Dataman-Cloud/seckilling/order/src/queue"
	_ "github.com/mattes/migrate/driver/mysql"
	migrate "github.com/mattes/migrate/migrate"
	"github.com/spf13/viper"
)

func main() {
	config.InitConfig()
	upgradeDB()
	queue.InitClient()
	queue.Start()
}

func upgradeDB() {
	user := viper.GetString("db.user")
	password := viper.GetString("db.password")
	host := viper.GetString("db.host")
	port := viper.GetInt("db.port")
	name := viper.GetString("db.name")
	driver := fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s", user, password, host, port, name)

	log.Println("upgrading DB", driver)
	errors, ok := migrate.UpSync(driver, "./sql")
	if errors != nil && len(errors) > 0 {
		for _, err := range errors {
			log.Panicln("db err", err)
		}
		log.Panicln("can't upgrade db", errors)
	}
	if !ok {
		log.Panicln("can't upgrade db")
	}
	log.Println("DB upgraded")
}
