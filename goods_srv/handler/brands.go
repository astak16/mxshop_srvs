package handler

import (
	"context"
	"mxshow_srvs/goods_srv/global"
	"mxshow_srvs/goods_srv/model"
	"mxshow_srvs/goods_srv/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 品牌和轮播图
func (g *GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	brandListResponse := proto.BrandListResponse{}

	var brands = []model.Brands{}

	// result := global.DB.Find(&brands)

	result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}

	var count int64
	global.DB.Model(&model.Brands{}).Count(&count)

	brandListResponse.Total = int32(count)

	var brandResponse []*proto.BrandInfoResponse
	for _, brand := range brands {
		brandResponse = append(brandResponse, &proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}
	brandListResponse.Data = brandResponse

	return &brandListResponse, nil
}

func (g *GoodsServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	if result := global.DB.First(&model.Brands{}); result.RowsAffected == 1 {
		return nil, status.Error(codes.InvalidArgument, "品牌已存在")
	}

	brand := &model.Brands{
		Name: req.Name,
		Logo: req.Logo,
	}

	global.DB.Save(brand)

	return &proto.BrandInfoResponse{Id: int32(brand.ID)}, nil
}
func (g *GoodsServer) DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*empty.Empty, error) {
	if result := global.DB.Delete(&model.Brands{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "品牌不存在")
	}
	return &empty.Empty{}, nil
}
func (g *GoodsServer) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*empty.Empty, error) {
	brands := model.Brands{}
	if result := global.DB.First(&brands); result.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "品牌不存在")
	}

	if req.Name != "" {
		brands.Name = req.Name
	}
	if req.Logo != "" {
		brands.Logo = req.Logo
	}

	global.DB.Save(&brands)
	return &empty.Empty{}, nil
}
