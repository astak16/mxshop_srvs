package main

import (
	"context"
	"fmt"
	"mxshow_srvs/goods_srv/proto"
)

func TestGoodsList() {
	rsp, err := barandClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
		TopCategory: 130361,
		// KeyWords:    "深海速冻",
		PriceMin: 90,
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, good := range rsp.Data {
		fmt.Println(good.Name, good.ShopPrice)
	}
}

func TestBatchGetGoods() {
	rsp, err := barandClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: []int32{421, 422, 423},
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, good := range rsp.Data {
		fmt.Println(good.Name, good.ShopPrice)
	}
}

func TestBatchGetGoodsDetail() {
	rsp, err := barandClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: 421,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(rsp.Name)
}
