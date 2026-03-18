package main

import (
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/srv/RabbitMQ"
	_ "github.com/Cao-yu-xin1/Spider_Crawl/lx0314/srv/basic/init1"
)

func main() {
	RabbitMQ.ConsumeStockDeduct()
}
