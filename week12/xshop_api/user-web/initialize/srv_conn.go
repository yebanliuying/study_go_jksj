package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important

	"xshop_api/user-web/global"
	"xshop_api/user-web/proto"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name), //&tag=srv
		//grpc.WithInsecure(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务】失败")
	}

	userClient := proto.NewUserClient(userConn)
	global.UserSrvCln = userClient
}

func InitSrvConnOld() {
	//从注册中心获取用户服务信息
	cfg := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)

	//本次获取数据的rpc服务器
	userSrvHost := ""
	userSrvPort := 0
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data,err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.ServerConfig.UserSrvInfo.Name))
	//data, err := client.Agent().ServicesWithFilter(`Service == "user-srv"`)

	if err != nil {
		panic(err)
	}

	//只获取一台服务器
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}

	if userSrvHost == "" {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务】失败")
	}

	//拨号连接用户grpc服务器
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务】 失败",
			"msg", err.Error(),
		)
	}

	//todo 1、后续的用户服务下线了 2、改端口了 3、该ip了 会出问题
	userClient := proto.NewUserClient(conn)
	global.UserSrvCln = userClient

}
