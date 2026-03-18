package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"lx0318/config"
)

var Rdb *redis.Client
var Ctx = context.Background()

func ExampleClient() {
	cng := config.GlobalConfig.Redis
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
