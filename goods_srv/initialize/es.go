package initialize

import (
	"context"
	"fmt"
	"log"
	"mxshow_srvs/goods_srv/global"
	"mxshow_srvs/goods_srv/model"
	"os"

	"github.com/olivere/elastic/v7"
)

func InitEs() {
	host := fmt.Sprintf("http://%s:%d", global.ServerConfig.EsInfo.Host, global.ServerConfig.EsInfo.Port)
	logger := log.New(os.Stdout, "mxshop", log.LstdFlags)
	// 这里必须将 sniff 设置为 false，因为使用 olivere/elastic 连接 es 时，发现连接地址是内网地址时，会自动转换成内网地址或者 docker 中的 ip 地址，导致连接不上
	var err error
	global.EsClient, err = elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false), elastic.SetTraceLog(logger))
	if err != nil {
		panic(err)
	}

	// 新建 mapping 和 index
	exists, err := global.EsClient.IndexExists(model.EsGoods{}.GetIndexName()).Do(context.Background())
	if err != nil {
		panic(err)
	}

	if !exists {
		_, err = global.EsClient.CreateIndex(model.EsGoods{}.GetIndexName()).BodyString(model.EsGoods{}.GetMapping()).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}
}
