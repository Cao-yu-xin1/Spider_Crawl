package cache

import (
	"context"
	"fmt"
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/nacos"
	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client
var Ctx = context.Background()

func ExampleClient() {
	cng := nacos.GlobalConfig.Redis
	Rdb = redis.NewClient(&redis.Options{
		Addr:     cng.Addr,
		Password: cng.Password, // no password set
		DB:       cng.Database, // use default DB
	})

	err := Rdb.Ping(Ctx).Err()
	if err != nil {
		panic(err)
	}
	fmt.Println("redis连接成功")
}
