package nacos

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
	viper.SetConfigFile("E:\\gowork\\src\\lx\\zg5\\lx0314\\nacos.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		panic("读取nacos配置文件失败" + err.Error())
	}
	var nacosConfig Nacos
	err = viper.UnmarshalKey("nacos", &nacosConfig)
	if err != nil {
		panic("nacos配置文件解析失败" + err.Error())
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
		panic("nacos创建配置客户端失败" + err.Error())
	}
	configContent, err := configClient.GetConfig(vo.ConfigParam{
		DataId: nacosConfig.DataId,
		Group:  nacosConfig.Group,
	})
	if err != nil {
		panic("nacos获取配置失败" + err.Error())
	}
	//fmt.Println(configContent)
	viper.SetConfigType("yaml")
	err = viper.ReadConfig(strings.NewReader(configContent))
	if err != nil {
		panic("nacos配置文件解析失败" + err.Error())
	}
	err = viper.Unmarshal(&GlobalConfig)
	if err != nil {
		panic("nacos配置文件解析失败" + err.Error())
	}
	//fmt.Println(GlobalConfig)
}
