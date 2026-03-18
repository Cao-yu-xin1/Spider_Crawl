package main

import (
	"context"
	"log"
	"lx0318/config"
	"lx0318/mq/RabbitMQ"
	"lx0318/pkg/cache"
	__ "lx0318/proto"
	"time"
)

func main() {
	kutengOne := RabbitMQ.NewRabbitMQTopic("exKutengTopic", "#")
	kutengOne.SubsribeMsg("topic", func(msg string) {
		val := cache.Rdb.SetNX(context.Background(), msg, 1, time.Minute*5).Val()
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
