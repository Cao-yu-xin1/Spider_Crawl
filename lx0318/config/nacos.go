package config

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"strings"
)

func App() {
	var err error
	viper.SetConfigFile("E:\\gowork\\src\\lx\\zg5\\lx0318\\nacos.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	var nacosConfig Nacos
	err = viper.UnmarshalKey("nacos", &nacosConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println(nacosConfig)
	// Nacos服务器地址
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: nacosConfig.Host,
			Port:   uint64(nacosConfig.Port),
		},
	}
	// 客户端配置
	clientConfig := constant.ClientConfig{
		NamespaceId:         nacosConfig.NamespaceId, // 如果不需要命名空间，可以留空
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	// 创建配置客户端
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}
	configContent, err := configClient.GetConfig(vo.ConfigParam{
		DataId: nacosConfig.DataId,
		Group:  nacosConfig.Group,
	})
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	err = viper.ReadConfig(strings.NewReader(configContent))
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&GlobalConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println(GlobalConfig)
}
