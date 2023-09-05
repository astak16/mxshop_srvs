package main

import (
	"flag"
	"fmt"
	"mxshow_srvs/goods_srv/global"
	"mxshow_srvs/goods_srv/handler"
	"mxshow_srvs/goods_srv/initialize"
	"mxshow_srvs/goods_srv/proto"
	"mxshow_srvs/goods_srv/utils"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50052, "端口号")

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	initialize.InitEs()
	zap.S().Info(global.ServerConfig)

	flag.Parse()

	zap.S().Info("ip: ", *IP)

	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	zap.S().Info("port: ", *Port)

	server := grpc.NewServer()
	proto.RegisterGoodsServer(server, &handler.GoodsServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	// 健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	check := &api.AgentServiceCheck{
		// CheckID:                        id,
		GRPC: fmt.Sprintf("%s:%d", global.ServerConfig.Host, *Port), //
		// HTTP:                           "http://" + address + ":" + strconv.Itoa(port) + "/health", // 检查服务的地址
		Timeout:                        "5s",   // 健康检查的超时时间
		Interval:                       "5s",   // 健康检查的间隔时间
		DeregisterCriticalServiceAfter: "400s", // 如果健康检查失败次数超过这个时间，则注销服务
	}
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	registration.ID = serviceID
	registration.Port = *Port
	registration.Tags = global.ServerConfig.Tags
	registration.Address = global.ServerConfig.Host
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc: " + err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("注销失败: ")
	}
	zap.S().Info("注销成功: ")
}
