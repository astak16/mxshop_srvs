package main

import (
	"flag"
	"fmt"
	"mxshow_srvs/order_srv/global"
	"mxshow_srvs/order_srv/handler"
	"mxshow_srvs/order_srv/initialize"
	"mxshow_srvs/order_srv/proto"
	"mxshow_srvs/order_srv/utils"
	"mxshow_srvs/order_srv/utils/register/consul"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50051, "端口号")

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	initialize.InitSrvConn()
	zap.S().Info(global.ServerConfig)

	flag.Parse()

	zap.S().Info("ip: ", *IP)

	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	zap.S().Info("port: ", *Port)

	server := grpc.NewServer()
	proto.RegisterOrderServer(server, &handler.OrderServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	// 健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 启动服务
	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc: " + err.Error())
		}
	}()

	fmt.Println(global.ServerConfig.Host, global.ServerConfig.ConsulInfo.Host)

	// 服务注册
	register_cliennt := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err = register_cliennt.Register(
		global.ServerConfig.Host,
		*Port,
		global.ServerConfig.Name,
		global.ServerConfig.Tags, serviceId,
	)
	if err != nil {
		zap.S().Panic("注册服务失败：", err.Error())
	}

	zap.S().Debugf("启动服务器，端口：%d", *Port)

	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"http://go-rmqnamesrv:9876"}),
		consumer.WithGroupName("mxshop-order"),
	)
	if err := c.Subscribe("order_timeout", consumer.MessageSelector{}, handler.AutoTimeout); err != nil {
		panic(err)
	}
	c.Start()

	// 终止退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	_ = c.Shutdown()

	if err = register_cliennt.DeRegister(serviceId); err != nil {
		zap.S().Info("注销服务失败：", err.Error())
	} else {
		zap.S().Info("注销服务成功")
	}
}
