package main

import (
	"context"
	"fmt"
	"mxshow_srvs/order_srv/proto"

	"google.golang.org/grpc"
)

var client proto.OrderClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client = proto.NewOrderClient(conn)
}

func TestCreateCartItem(userId, nums, goodsId int32) {
	rsp, err := client.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  userId,
		Nums:    nums,
		GoodsId: goodsId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id)
}

func TestCartItemList(userId int32) {
	rsp, err := client.CartItemList(context.Background(), &proto.UserInfo{
		Id: userId,
	})
	if err != nil {
		panic(err)
	}
	for _, item := range rsp.Data {
		fmt.Println(item.Id, item.GoodsId, item.Nums)
	}
}

func TestUpdateCartItem(id int32) {
	_, err := client.UpdateCartItem(context.Background(), &proto.CartItemRequest{
		Id:      id,
		Checked: false,
	})
	if err != nil {
		panic(err)
	}
}

func TestCreateOrder() {
	_, err := client.Create(context.Background(), &proto.OrderRequest{
		UserId:  1,
		Address: "北京",
		Name:    "boby",
		Mobile:  "123456789",
		Post:    "请尽快发货",
	})
	if err != nil {
		panic(err)
	}
}

func TestGetOrderDetail(orderId int32) {
	rsp, err := client.OrderDetail(context.Background(), &proto.OrderRequest{
		Id: orderId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.OrderInfo.OrderSn)
	for _, good := range rsp.Goods {
		fmt.Println(good.GoodsName)
	}
}

func TestOrderList() {
	rsp, err := client.OrderList(context.Background(), &proto.OrderFilterRequest{
		UserId: 1,
	})
	if err != nil {
		panic(err)
	}

	for _, order := range rsp.Data {
		fmt.Println(order.OrderSn)
	}
}

func main() {
	Init()
	defer conn.Close()
	// TestCreateCartItem(2, 1, 421)
	// TestCartItemList(1)
	// TestUpdateCartItem(2)
	// TestCreateOrder()
	// TestGetOrderDetail(2)
	TestOrderList()
}
