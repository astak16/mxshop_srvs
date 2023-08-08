package main

import (
	"context"
	"fmt"
	"mxshow_srvs/goods_srv/proto"

	"github.com/golang/protobuf/ptypes/empty"
)

func TestGetCategoryList() {
	rsp, err := barandClient.GetAllCategorysList(context.Background(), &empty.Empty{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	fmt.Println(rsp.JsonData)
	// for _, brand := range rsp.Data {
	// 	fmt.Println(brand.Name)

	// }
}

func TestGetSubCategoryList() {
	rsp, err := barandClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id: 130358,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.SubCategorys)
	// fmt.Println(rsp.JsonData)
	// for _, brand := range rsp.Data {
	// 	fmt.Println(brand.Name)

	// }
}
