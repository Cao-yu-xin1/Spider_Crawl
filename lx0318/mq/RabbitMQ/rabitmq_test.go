package RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"testing"
)

func TestRabbitMQ_SendMsg(t *testing.T) {
	conn, _ := amqp.Dial("amqp://caoyuxin:caoyuxin@115.190.154.22:5672/xin")
	channel, _ := conn.Channel()
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
				conn:      conn,
				channel:   channel,
				QueueName: "1",
				Exchange:  "2",
				Key:       "3",
				Mqurl:     "amqp://caoyuxin:caoyuxin@115.190.154.22:5672/xin",
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
	conn, _ := amqp.Dial("amqp://caoyuxin:caoyuxin@115.190.154.22:5672/xin")
	channel, _ := conn.Channel()
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
				conn:      conn,
				channel:   channel,
				QueueName: "1",
				Exchange:  "4",
				Key:       "3",
				Mqurl:     "amqp://caoyuxin:caoyuxin@115.190.154.22:5672/xin",
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
