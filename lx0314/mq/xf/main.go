package main

import (
	"context"
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/bff-api/basic/config"
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/mq/rabbit"
	__ "github.com/Cao-yu-xin1/Spider_Crawl/lx0314/proto"
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/srv/pkg/cache"
	"log"
	"time"
)

func main() {
	kutengOne := rabbit.NewRabbitMQTopic("exKutengTopic", "#")
	kutengOne.RecieveTopic("topic", func(msg string) {
		val := cache.Rdb.SetNX(context.Background(), "order_sn", 1, time.Minute*5).Val()
		if !val {
			log.Println("订单已处理或正在处理中")
			return
		}
		log.Printf("收到消息: %s", msg)

		// 调用通知支付服务
		_, err := config.ServiceClient.NotifyPay(
			context.Background(),
			&__.NotifyPayReq{OrderNo: msg},
		)
		if err != nil {
			log.Printf("通知支付失败: %v", err)
			// 这里可以考虑重试或记录失败
		}
	})
}
