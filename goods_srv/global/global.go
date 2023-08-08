package global

import (
	"mxshow_srvs/goods_srv/config"

	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
)

// docker run -d --network network1 --network-alias go-consul  -p 8500:8500 -p 8300:8300 -p 8301:8301 -p 8302:8302 -p 8600:8600/udp --name go-consul  consul:1.15 consul agent -dev -client=0.0.0.0
// docker run -d --network network1 --network-alias redis -p 6379:6379 --name go-redis redis:latest
