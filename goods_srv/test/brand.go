package main

import (
	"context"
	"fmt"
	"mxshow_srvs/goods_srv/proto"

	"google.golang.org/grpc"
)

var barandClient proto.GoodsClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	barandClient = proto.NewGoodsClient(conn)
}

func TestGetBrandList() {
	rsp, err := barandClient.BrandList(context.Background(), &proto.BrandFilterRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, brand := range rsp.Data {
		fmt.Println(brand.Name)

	}
}
