package publish

import (
	"lx0318/mq/RabbitMQ"
)

func SendMessage(topic string, orderNo string) {
	kutengOne := RabbitMQ.NewRabbitMQTopic("exKutengTopic", "kuteng.topic.one")
	kutengOne.SendMsg(topic, orderNo)
}
