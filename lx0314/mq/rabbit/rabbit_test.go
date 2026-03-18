package rabbit

import (
	"github.com/streadway/amqp"
	"testing"
)

func TestRabbitMQ_PublishTopic(t *testing.T) {
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
		message string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "发布到默认topic",
			fields: fields{
				conn:      nil, // 需要连接实际RabbitMQ或使用模拟
				channel:   nil, // 需要实际通道或模拟
				QueueName: "test_queue",
				Exchange:  "test_exchange",
				Key:       "test_key",
				Mqurl:     "amqp://guest:guest@localhost:5672/",
			},
			args: args{
				topic:   "test.topic",
				message: "Hello World",
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
			r.PublishTopic(tt.args.topic, tt.args.message)
		})
	}
}
