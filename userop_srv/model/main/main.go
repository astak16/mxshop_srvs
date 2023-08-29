package main

import (
	"fmt"
	"log"
	"mxshow_srvs/userop_srv/model"
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
	mysqlDbname   = "mxshop_userop_srv"
)

// func genMd5(code string) string {
// 	Md5 := md5.New()
// 	_, _ = io.WriteString(Md5, code)
// 	return hex.EncodeToString(Md5.Sum(nil))
// }

func main() {
	dsn := fmt.Sprintf("%s:%d@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlUser, mysqlPassword, mysqlHost, mysqlDbname)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable colors
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})

	// fmt.Println(&db, err)

	if err != nil {
		panic(err)
	}
	_ = db.AutoMigrate(&model.Address{}, &model.LeavingMessages{}, &model.UserFav{})
}
