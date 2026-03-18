package init1

import (
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/nacos"
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/srv/pkg/cache"
)

func init() {
	nacos.App()
	MysqlInit()
	cache.ExampleClient()
	ConsulInit()
}
