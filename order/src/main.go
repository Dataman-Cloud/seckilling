package main

import (
	"github.com/Dataman-Cloud/seckilling/order/src/config"
	"github.com/Dataman-Cloud/seckilling/order/src/queue"
)

func main() {
	config.InitConfig()
	queue.InitClient()
	queue.Start()
}
