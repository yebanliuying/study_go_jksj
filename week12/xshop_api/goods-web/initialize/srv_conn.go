package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"xshop_api/goods-web/global"
	"xshop_api/goods-web/proto"
	"xshop_api/goods-web/utils/otgrpc"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name), //&tag=srv
		//grpc.WithInsecure(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())), //链路追踪
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【商品服务】失败")
	}

	goodsClient := proto.NewGoodsClient(userConn)
	global.GoodsSrvCln = goodsClient
}
