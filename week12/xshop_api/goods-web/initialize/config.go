package initialize

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"xshop_api/goods-web/global"
)

//GetEnvInfo 获取环境变量
func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	dev := GetEnvInfo("XSHOP_DEV")

	if dev {
		global.ServerConfig.Env = "dev"
	} else {
		global.ServerConfig.Env = "pro"
	}



	configFilePrefix := "config"
	configFileName := fmt.Sprintf("goods-web/%s-%s.yaml",configFilePrefix,global.ServerConfig.Env)


	//设置文件路径
	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	//获取config server配置
	//serverConfig := config.ServerConfig{}
	if err := v.Unmarshal(global.NacosConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息: %v", global.NacosConfig)


	//从nacos中读取配置信息
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,//配置中心
			Port: global.NacosConfig.Port,
		},
	}

	cc := constant.ClientConfig {
		NamespaceId:         global.NacosConfig.Namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	// 创建动态配置客户端
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}

	//获取配置信息
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group})

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("读取nacos配置失败:%s", err.Error())
	}
	fmt.Println(&global.ServerConfig)
}