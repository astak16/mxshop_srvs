package main

import (
	"context"
	"fmt"
	"mxshow_srvs/goods_srv/proto"
)

func TestGetCategoryBrandList() {
	rsp, err := barandClient.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Data)
	// fmt.Println(rsp.JsonData)
	// for _, brand := range rsp.Data {
	// 	fmt.Println(brand.Name)

	// }
}

// func TestGetSubCategoryList() {
// 	rsp, err := barandClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
// 		Id: 130358,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(rsp.SubCategorys)
// 	// fmt.Println(rsp.JsonData)
// 	// for _, brand := range rsp.Data {
// 	// 	fmt.Println(brand.Name)

// 	// }
// }
