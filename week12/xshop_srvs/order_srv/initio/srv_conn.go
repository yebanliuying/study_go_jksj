package initio

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"xshop_srvs/order_srv/global"
	"xshop_srvs/order_srv/proto"
)

func InitSrvConn() {
	//初始化商品服务连接
	consulInfo := global.ServerConfig.ConsulInfo
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name), //&tag=srv
		//grpc.WithInsecure(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【商品服务】失败")
	}

	goodsClient := proto.NewGoodsClient(goodsConn)
	global.GoodsSrvClient = goodsClient

	//初始化库存服务连接
	//consulInfo := global.ServerConfig.ConsulInfo
	inventoryConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.InventorySrvInfo.Name), //&tag=srv
		//grpc.WithInsecure(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【库存服务】失败")
	}

	inventoryClient := proto.NewInventoryClient(inventoryConn)
	global.InventorySrvClient = inventoryClient
}
