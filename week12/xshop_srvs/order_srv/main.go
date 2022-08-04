package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

	"xshop_srvs/order_srv/global"
	"xshop_srvs/order_srv/handler"
	"xshop_srvs/order_srv/initio"
	"xshop_srvs/order_srv/proto"
	"xshop_srvs/order_srv/utils"
	"xshop_srvs/order_srv/utils/register/consul"
	"xshop_srvs/order_srv/utils/otgrpc"
)

func main() {

	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50051, "端口号")

	//初始化
	initio.InitLogger()
	initio.InitConfig()
	initio.InitDB()
	initio.InitRS()
	initio.InitSrvConn() //第三方微服务

	flag.Parse()
	zap.S().Info("ip:", *IP)
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}
	zap.S().Info("port:", *Port)


	//生成一个新的tracer
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: fmt.Sprintf("%s:%d", global.ServerConfig.JaegerInfo.Host, global.ServerConfig.JaegerInfo.Port),
		},
		ServiceName: global.ServerConfig.JaegerInfo.Name,
	}
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}

	//开启全局tracer
	opentracing.SetGlobalTracer(tracer)
	//开启服务 接受otgrpc 链路追踪
	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))

	//注册服务
	//proto.RegisterOrderServer(server, &proto.UnimplementedOrderServer{})
	proto.RegisterOrderServer(server, &handler.OrderServer{})

	//监听端口
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
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


	//监听订单超时topic
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"139.198.21.42:9876"}),
		consumer.WithGroupName("xshop-order"),
	)

	if err := c.Subscribe("order_timeout",  consumer.MessageSelector{}, handler.OrderTimeout ); err != nil {
		fmt.Println("读取消息失败")
	}

	_ = c.Start()
	//不能让主goroutine 推出
	//time.Sleep(time.Hour)



	//接受终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	_ = c.Shutdown() //一个进程中只能有一个shutdown 除非是创建多个consumer
	_ = closer.Close()
	if err = registerClient.DeRegister(serviceId); err != nil {
		zap.S().Panic("注销失败:", err.Error())
	} else {
		zap.S().Info("注销成功")
	}

}
