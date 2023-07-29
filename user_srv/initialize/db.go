package initialize

import (
	"fmt"
	"log"
	"mxshow_srvs/user_srv/global"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// const (
// 	mysqlHost     = "go-uccs-1"
// 	mysqlPort     = 3306
// 	mysqlUser     = "root"
// 	mysqlPassword = 123456
// 	mysqlDbname   = "mxshop_user_srv"
// )

func InitDB() {
	c := global.ServerConfig.MysqlInfo
	// dsn := fmt.Sprintf("%s:%d@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlUser, mysqlPassword, mysqlHost, mysqlDbname)
	dsn := fmt.Sprintf("%s:%d@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.Name)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Log level
			Colorful:      true,          // Disable colors
		},
	)

	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}
}
