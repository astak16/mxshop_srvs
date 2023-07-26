package global

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	mysqlHost     = "go-uccs-1"
	mysqlPort     = 3306
	mysqlUser     = "root"
	mysqlPassword = 123456
	mysqlDbname   = "mxshop_user_srv"
)

var (
	DB *gorm.DB
)

func init() {
	dsn := fmt.Sprintf("%s:%d@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlUser, mysqlPassword, mysqlHost, mysqlDbname)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Log level
			Colorful:      false,         // Disable colors
		},
	)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}
}

// docker run -d --network network1 --network-alias go-consul  -p 8500:8500 -p 8300:8300 -p 8301:8301 -p 8302:8302 -p 8600:8600/udp --name go-consul  consul:1.15 consul agent -dev -client=0.0.0.0
// docker run -d --network network1 --network-alias redis -p 6379:6379 --name go-redis redis:latest
