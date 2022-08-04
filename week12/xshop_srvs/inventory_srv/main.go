package main

import (
	"flag"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"net"
	"os"
	"os/signal"
	"syscall"
	"xshop_srvs/inventory_srv/handler"
	"xshop_srvs/inventory_srv/utils/register/consul"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"xshop_srvs/inventory_srv/global"
	"xshop_srvs/inventory_srv/initio"
	"xshop_srvs/inventory_srv/proto"
	"xshop_srvs/inventory_srv/utils"
)

func main() {

	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50053, "端口号")

	//初始化
	initio.InitLogger()
	initio.InitConfig()
	initio.InitDB()
	initio.InitRS()

	flag.Parse()
	zap.S().Info("ip:", *IP)
	if *Port == 0{
		*Port, _ = utils.GetFreePort()
	}
	zap.S().Info("port:", *Port)

	//开启服务
	server := grpc.NewServer()

	//注册服务
	proto.RegisterInventoryServer(server, &handler.InventoryServer{})

	//监听端口
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host,global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err = registerClient.Register(global.ServerConfig.Host, *Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败:", err.Error())
	}
	zap.S().Infof("启动服务器，端口：%d", *Port)

	//启动服务
	go func() {
		//异步调用堵塞服务
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()


	//监听库存归还topic
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"139.198.21.42:9876"}),
		consumer.WithGroupName("group1"),
	)

	if err := c.Subscribe("order_reback",  consumer.MessageSelector{}, handler.AutoReback ); err != nil {
		fmt.Println("读取消息失败")
	}

	_ = c.Start()
	//不能让主goroutine 推出
	//time.Sleep(time.Hour)




	//接受终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	_ = c.Shutdown()
	if err = registerClient.DeRegister(serviceId); err != nil {
		zap.S().Panic("注销失败:", err.Error())
	}else{
		zap.S().Info("注销成功")
	}

}
