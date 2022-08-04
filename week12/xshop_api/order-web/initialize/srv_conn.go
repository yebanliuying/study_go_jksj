package initialize

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important

	"xshop_api/order-web/global"
	"xshop_api/order-web/proto"
	"xshop_api/order-web/utils/otgrpc"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	orderConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",consulInfo.Host, consulInfo.Port, global.ServerConfig.OrderSrvInfo.Name), //&tag=srv
		//grpc.WithInsecure(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())), //链路追踪
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【订单服务】失败")
	}

	orderClient := proto.NewOrderClient(orderConn)
	global.OrderSrvCln = orderClient

	//consulInfo := global.ServerConfig.ConsulInfo
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name), //&tag=srv
		//grpc.WithInsecure(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())), //链路追踪
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【商品服务】失败")
	}

	goodsClient := proto.NewGoodsClient(goodsConn)
	global.GoodsSrvCln = goodsClient

	//consulInfo := global.ServerConfig.ConsulInfo
	inventoryConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",consulInfo.Host, consulInfo.Port, global.ServerConfig.InventorySrvInfo.Name), //&tag=srv
		//grpc.WithInsecure(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())), //链路追踪
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【库存服务】失败")
	}

	inventoryClient := proto.NewInventoryClient(inventoryConn)
	global.InventorySrvCln = inventoryClient
}
