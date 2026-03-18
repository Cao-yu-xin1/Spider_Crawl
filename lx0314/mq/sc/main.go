package sc

import "github.com/Cao-yu-xin1/Spider_Crawl/lx0314/mq/rabbit"

func SendMessage(topic string, msg string) {
	kutengOne := rabbit.NewRabbitMQTopic("exKutengTopic", "kuteng.topic.one")
	kutengOne.PublishTopic(topic, msg)
}
