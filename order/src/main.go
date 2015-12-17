package main

import (
	"log"

	"github.com/Dataman-Cloud/seckilling/order/src/cache"
	"github.com/Dataman-Cloud/seckilling/order/src/config"
	"github.com/Dataman-Cloud/seckilling/order/src/db"
	"github.com/Dataman-Cloud/seckilling/order/src/queue"
)

func componentInit() {
	config.InitConfig()
	cache.InitCache()
	db.InitDB()
}

func destroy() {
	log.Println("destroy...")
	cache.DestroyCache()
	db.DestroyDB()
}

func main() {
	componentInit()
	defer destroy()

	db.UpgradeDB()
	queue.InitClient()
	queue.Start()
}
