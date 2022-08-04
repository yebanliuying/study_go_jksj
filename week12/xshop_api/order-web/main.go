package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"os"
	"os/signal"
	"syscall"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"

	"xshop_api/order-web/global"
	"xshop_api/order-web/initialize"
	"xshop_api/order-web/utils"
	"xshop_api/order-web/utils/register/consul"
	"xshop_api/order-web/validate"
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

	//开发环境固定端口号，线上环境自动获取端口号
	if global.ServerConfig.Env == "pro" {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}

	}


	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", validate.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
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
