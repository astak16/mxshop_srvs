package handler

import (
	"context"
	"encoding/json"
	"mxshow_srvs/goods_srv/global"
	"mxshow_srvs/goods_srv/model"
	"mxshow_srvs/goods_srv/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// 商品分类
func (g *GoodsServer) GetAllCategorysList(ctx context.Context, req *empty.Empty) (*proto.CategoryListResponse, error) {
	var categorys = []model.Category{}
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)

	b, _ := json.Marshal(categorys)

	return &proto.CategoryListResponse{JsonData: string(b)}, nil
}

// 获取子分类
func (g *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	categoryListResponse := proto.SubCategoryListResponse{}

	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	categoryListResponse.Info = &proto.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		Level:          category.Level,
		IsTab:          category.IsTab,
		ParentCategory: category.ParentCategoryId,
	}

	var subCategorys []model.Category
	var subCategoryResponse []*proto.CategoryInfoResponse
	preloads := "SubCategory"
	if category.Level == 1 {
		preloads = "SubCategory.SubCategory"
	}

	global.DB.Where(&model.Category{ParentCategoryId: category.ID}).Preload(preloads).Find(&subCategorys)

	for _, subCategory := range subCategorys {
		subCategoryResponse = append(subCategoryResponse, &proto.CategoryInfoResponse{
			Id:             subCategory.ID,
			Name:           subCategory.Name,
			Level:          subCategory.Level,
			IsTab:          subCategory.IsTab,
			ParentCategory: subCategory.ParentCategoryId,
		})
	}

	categoryListResponse.SubCategorys = subCategoryResponse

	return &categoryListResponse, nil
}
func (s *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	category := model.Category{}

	category.Name = req.Name
	category.Level = req.Level
	if req.Level != 1 {
		category.ParentCategoryId = req.ParentCategory
	}
	category.IsTab = req.IsTab

	global.DB.Save(&category)

	return &proto.CategoryInfoResponse{Id: int32(category.ID)}, nil
}

func (s *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Category{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	var category model.Category

	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentCategory != 0 {
		category.ParentCategoryId = req.ParentCategory
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	if req.IsTab {
		category.IsTab = req.IsTab
	}

	global.DB.Save(&category)

	return &emptypb.Empty{}, nil
}
