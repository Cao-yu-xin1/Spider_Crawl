package init

import (
	"lx0318/config"
	"lx0318/pkg/cache"
)

func init() {
	config.App()
	MysqlInit()
	cache.ExampleClient()
}
