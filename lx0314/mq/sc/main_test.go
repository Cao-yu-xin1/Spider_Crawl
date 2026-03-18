package sc

import "testing"

func TestSendMessage(t *testing.T) {
	type args struct {
		topic string
		msg   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "测试1",
			args: args{
				topic: "topic",
				msg:   "bbb",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendMessage(tt.args.topic, tt.args.msg)
		})
	}
}
