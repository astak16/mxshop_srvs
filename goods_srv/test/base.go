package main

func main() {
	Init()
	defer conn.Close()
	// TestGetBrandList()
	// TestGetCategoryList()
	// TestGetSubCategoryList()
	// TestGetCategoryBrandList()
	// TestGoodsList()
	// TestBatchGetGoods()
	TestBatchGetGoodsDetail()
}
