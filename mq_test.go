package mqgo

import (
	"fmt"
	"testing"
	"time"
)

var (
	mqConnStr    = "amqp://admin:admin@10.1.1.248:5672/"
	mysqlConnStr = "root:pwd@(10.1.1.245:3306)/datacenter?charset=utf8"
	queue        = "dc-sz-test-mqsync"
	key          = "dc-sz-test-mqsync"
	exchange     = "datacenter"
)

func TestMq_Publish(t *testing.T) {
	type args struct {
		model SyncMqInfo
	}
	tests := []struct {
		name    string
		mq      *Mq
		args    args
		wantErr bool
	}{
		{
			name: "订阅",
			mq:   NewMqByStr(mqConnStr, mysqlConnStr),
			args: args{
				model: SyncMqInfo{
					Exchange: exchange,
					Queue:    queue,
					RouteKey: key,
					Request:  time.Now().String(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.mq.Publish(tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("Mq.Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMq_Consume(t *testing.T) {
	type args struct {
		queue    string
		key      string
		exchange string
		fun      func(request string) (response string, err error)
	}
	tests := []struct {
		name string
		mq   *Mq
		args args
	}{
		{
			name: "消费",
			mq:   NewMqByStr(mqConnStr, mysqlConnStr),
			args: args{
				queue:    queue,
				key:      key,
				exchange: exchange,
				fun: func(request string) (string, error) {
					fmt.Println(request)
					return "测试成功", nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mq.Consume(tt.args.queue, tt.args.key, tt.args.exchange, tt.args.fun)
		})
	}
}
