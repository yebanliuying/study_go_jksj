package initialize

import "go.uber.org/zap"

func InitLogger()  {

	//logger, _ := zap.NewProduction() //todo 根据配置的dev和prd来控制log类型,配置项再暴露修改log输出类型
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}