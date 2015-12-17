package db

import (
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattes/migrate/driver/mysql"
	migrate "github.com/mattes/migrate/migrate"
	"github.com/spf13/viper"
)

var db *sqlx.DB

func DB() *sqlx.DB {
	if db != nil {
		return db
	}

	mutex := sync.Mutex{}
	mutex.Lock()
	InitDB()
	defer mutex.Unlock()

	return db
}

func InitDB() {
	dburl := DBUrl()
	url := fmt.Sprintf("%s?parseTime=true", dburl)
	var err error
	db, err = sqlx.Open("mysql", url)
	if err != nil {
		log.Panicln("can't open db: ", url, " err: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Panicln("can't ping db: ", url, " err: ", err)
	}

	maxIdleConns := viper.GetInt("db.maxIdleConns")
	maxOpenConns := viper.GetInt("db.maxOpenConns")

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)

	log.Println("initialized db: ", url)
}

func DestroyDB() {
	log.Println("destroying DB")
	if db != nil {
		db.Close()
		log.Println("db closed")
	}
}

func UpgradeDB() {
	dburl := DBUrl()
	driver := fmt.Sprintf("mysql://%s", dburl)

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

func DBUrl() string {
	user := viper.GetString("db.user")
	password := viper.GetString("db.password")
	host := viper.GetString("db.host")
	port := viper.GetInt("db.port")
	name := viper.GetString("db.name")
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, name)
}
