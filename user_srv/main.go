package main

import (
	"flag"
	"fmt"
	"mxshow_srvs/user_srv/global"
	"mxshow_srvs/user_srv/handler"
	"mxshow_srvs/user_srv/initialize"
	"mxshow_srvs/user_srv/proto"
	"net"

	"github.com/hashicorp/consul/api"
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
	zap.S().Info("port: ", *Port)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})

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
		GRPC: fmt.Sprintf("172.18.0.3:50051"), // 检查服务的 ID
		// HTTP:                           "http://" + address + ":" + strconv.Itoa(port) + "/health", // 检查服务的地址
		Timeout:                        "5s",   // 健康检查的超时时间
		Interval:                       "5s",   // 健康检查的间隔时间
		DeregisterCriticalServiceAfter: "400s", // 如果健康检查失败次数超过这个时间，则注销服务
	}
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	registration.ID = global.ServerConfig.Name
	registration.Port = *Port
	registration.Tags = []string{"imooc", "bobby", "user", "srv"}
	registration.Address = "172.18.0.3"
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	err = server.Serve(lis)
	if err != nil {
		panic("failed to start grpc: " + err.Error())
	}
}
