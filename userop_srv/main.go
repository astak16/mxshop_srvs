package main

import (
	"flag"
	"fmt"

	"mxshow_srvs/userop_srv/global"
	"mxshow_srvs/userop_srv/handler"
	"mxshow_srvs/userop_srv/initialize"
	"mxshow_srvs/userop_srv/proto"
	"mxshow_srvs/userop_srv/utils"
	"mxshow_srvs/userop_srv/utils/register/consul"
	"net"
	"os"
	"os/signal"
	"syscall"

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
	zap.S().Info(global.ServerConfig)

	flag.Parse()

	zap.S().Info("ip: ", *IP)

	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	zap.S().Info("port: ", *Port)

	server := grpc.NewServer()
	proto.RegisterMessageServer(server, &handler.UserOpServer{})
	proto.RegisterAddressServer(server, &handler.UserOpServer{})
	proto.RegisterUserFavServer(server, &handler.UserOpServer{})

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

	// 终止退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = register_cliennt.DeRegister(serviceId); err != nil {
		zap.S().Info("注销服务失败：", err.Error())
	} else {
		zap.S().Info("注销服务成功")
	}
}
