package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mxshow_srvs/goods_srv/global"
	"mxshow_srvs/goods_srv/model"
	"os"
	"strconv"
	"time"

	"github.com/olivere/elastic/v7"
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
	mysqlDbname   = "mxshop_goods_srv"
	esHost        = "go-elasticsearch"
	esPort        = 9200
)

func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}

func main() {
	// dsn := fmt.Sprintf("%s:%d@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlUser, mysqlPassword, mysqlHost, mysqlDbname)
	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags),
	// 	logger.Config{
	// 		SlowThreshold: time.Second,   // Slow SQL threshold
	// 		LogLevel:      logger.Silent, // Log level
	// 		Colorful:      false,         // Disable colors
	// 	},
	// )

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	// 	NamingStrategy: schema.NamingStrategy{
	// 		SingularTable: true,
	// 	},
	// 	Logger: newLogger,
	// })

	// if err != nil {
	// 	panic(err)
	// }
	// _ = db.AutoMigrate(
	// 	&model.Category{},
	// 	&model.Brands{},
	// 	&model.GoodsCategoryBrand{},
	// 	&model.Banner{},
	// 	&model.Goods{},
	// )
	Mysql2Es()
}

func Mysql2Es() {
	// mysql 初始化
	dsn := fmt.Sprintf("%s:%d@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlUser, mysqlPassword, mysqlHost, mysqlDbname)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Log level
			Colorful:      false,         // Disable colors
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	// es 初始化
	host := fmt.Sprintf("http://%s:%d", esHost, esPort)
	logger := log.New(os.Stdout, "mxshop", log.LstdFlags)
	// 这里必须将 sniff 设置为 false，因为使用 olivere/elastic 连接 es 时，发现连接地址是内网地址时，会自动转换成内网地址或者 docker 中的 ip 地址，导致连接不上
	global.EsClient, err = elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false), elastic.SetTraceLog(logger))
	if err != nil {
		panic(err)
	}

	// 同步 goods
	var goods []model.Goods
	db.Find(&goods)
	for _, good := range goods {
		esModel := model.EsGoods{
			ID:          good.ID,
			CategoryID:  good.CategoryId,
			BrandsID:    good.BrandsId,
			OnSale:      good.OnSale,
			ShipFree:    good.ShipFree,
			IsNew:       good.IsNew,
			IsHot:       good.IsHot,
			Name:        good.Name,
			ClickNum:    good.ClickNum,
			SoldNum:     good.SoldNum,
			FavNum:      good.FavNum,
			MarketPrice: good.MarketPrice,
			ShopPrice:   good.ShopPrice,
			GoodsBrief:  good.GoodsBrief,
		}
		// 第一个 Index 是动词，表示要执行的操作；第二个 Index 是名词，表示要操作的索引
		_, err := global.EsClient.Index().Index(esModel.GetIndexName()).BodyJson(esModel).Id(strconv.Itoa(int(good.ID))).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}

}
