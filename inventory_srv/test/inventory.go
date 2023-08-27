package main

import (
	"context"
	"fmt"
	"mxshow_srvs/inventory_srv/proto"

	"google.golang.org/grpc"
)

var invClient proto.InventoryClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	invClient = proto.NewInventoryClient(conn)
}

func TestSetInv(goodsId, Num int32) {
	_, err := invClient.SetInv(context.Background(), &proto.GoodsInvInfo{
		GoodsId: goodsId,
		Num:     Num,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("设置库存成功")
}

func TestInvDetail(goodsId int32) {
	rsp, err := invClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: goodsId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Num)
}

func TestSell() {
	_, err := invClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId: 421, Num: 10},
			{GoodsId: 422, Num: 50},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("库存扣减成功")
}

func TestReback() {
	_, err := invClient.Reback(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId: 421, Num: 10},
			{GoodsId: 422, Num: 50},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("归还成功")
}

func main() {
	Init()
	defer conn.Close()
	var i int32
	for i = 421; i < 840; i++ {
		TestSetInv(i, 100)
	}
	// TestSetInv(422, 70)
	// TestInvDetail(421)
	// TestSell()
	// TestReback()
}
