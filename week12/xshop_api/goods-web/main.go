package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"

	"xshop_api/goods-web/global"
	"xshop_api/goods-web/initialize"
	"xshop_api/goods-web/utils"
	"xshop_api/goods-web/utils/register/consul"
)



func main() {

	//初始化log
	initialize.InitLogger()

	//初始化配置文件
	initialize.InitConfig()

	//初始化路由
	Router := initialize.InitRouters()

	//初始化翻译器
	initialize.InitTrans("zh")

	//初始化srv的连接
	initialize.InitSrvConn()

	//初始化sentinel 处理限流、熔断
	initialize.InitSentinel()

	//开发环境固定端口号，线上环境自动获取端口号
	if global.ServerConfig.Env == "pro" {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}

	}

	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host,global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err := registerClient.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败:", err.Error())
	}
	zap.S().Infof("启动服务器，端口：%d", global.ServerConfig.Port)

	go func(){
		if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动失败:", err.Error())
		}
	}()

	//接受终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = registerClient.DeRegister(serviceId); err != nil {
		zap.S().Panic("注销失败:", err.Error())
	}else{
		zap.S().Info("注销成功")
	}


}
