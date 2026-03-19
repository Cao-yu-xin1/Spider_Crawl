package RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"testing"
)

func TestRabbitMQ_SendMsg(t *testing.T) {
	topic := NewRabbitMQTopic("exKutengTopic", "kuteng.topic.one")
	//conn, _ := amqp.Dial("amqp://caoyuxin:caoyuxin@115.190.154.22:5672/xin")
	//channel, _ := conn.Channel()
	defer topic.channel.Close()
	type fields struct {
		conn      *amqp.Connection
		channel   *amqp.Channel
		QueueName string
		Exchange  string
		Key       string
		Mqurl     string
	}
	type args struct {
		topic string
		msg   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "1",
			fields: fields{
				conn:      topic.conn,
				channel:   topic.channel,
				QueueName: topic.QueueName,
				Exchange:  topic.Exchange,
				Key:       topic.Key,
				Mqurl:     topic.Mqurl,
			},
			args: args{
				topic: "topic",
				msg:   "1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RabbitMQ{
				conn:      tt.fields.conn,
				channel:   tt.fields.channel,
				QueueName: tt.fields.QueueName,
				Exchange:  tt.fields.Exchange,
				Key:       tt.fields.Key,
				Mqurl:     tt.fields.Mqurl,
			}
			r.SendMsg(tt.args.topic, tt.args.msg)
		})
	}
}

func TestRabbitMQ_SubsribeMsg(t *testing.T) {
	topic := NewRabbitMQTopic("exKutengTopic", "kuteng.topic.one")
	defer topic.channel.Close()
	//conn, _ := amqp.Dial("amqp://caoyuxin:caoyuxin@115.190.154.22:5672/xin")
	//channel, _ := conn.Channel()
	type fields struct {
		conn      *amqp.Connection
		channel   *amqp.Channel
		QueueName string
		Exchange  string
		Key       string
		Mqurl     string
	}
	type args struct {
		topic   string
		handler func(msg string)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "1",
			fields: fields{
				conn:      topic.conn,
				channel:   topic.channel,
				QueueName: topic.QueueName,
				Exchange:  topic.Exchange,
				Key:       topic.Key,
				Mqurl:     topic.Mqurl,
			},
			args: args{
				topic: "topic",
				handler: func(msg string) {
					fmt.Println(msg)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RabbitMQ{
				conn:      tt.fields.conn,
				channel:   tt.fields.channel,
				QueueName: tt.fields.QueueName,
				Exchange:  tt.fields.Exchange,
				Key:       tt.fields.Key,
				Mqurl:     tt.fields.Mqurl,
			}
			r.SubsribeMsg(tt.args.topic, tt.args.handler)
		})
	}
}
