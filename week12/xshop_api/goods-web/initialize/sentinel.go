package initialize

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"go.uber.org/zap"
)

func InitSentinel() {
	//先初始化sentinel
	err := sentinel.InitDefault()
	if err != nil {
		zap.S().Fatalf("初始化sentinel 异常:%v", err)
	}

	//配置限流规则
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "goods-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,//直接拒绝
			//ControlBehavior:        flow.Throttling,//匀速通过
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
	})

	if err != nil {
		zap.S().Fatalf("加载规则失败：%v", err)
	}
}