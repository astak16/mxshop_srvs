package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"mxshow_srvs/goods_srv/global"
	"mxshow_srvs/goods_srv/model"
	"mxshow_srvs/goods_srv/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/olivere/elastic/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

func ModelToResponse(goods model.Goods) proto.GoodsInfoResponse {
	return proto.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryId,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brands.ID,
			Name: goods.Brands.Name,
			Logo: goods.Brands.Logo,
		},
	}

}

// 商品接口
func (s *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	goodsListResponse := &proto.GoodsListResponse{}
	localDB := global.DB.Model(model.Goods{})

	q := elastic.NewBoolQuery()

	if req.KeyWords != "" {
		q.Must(elastic.NewMultiMatchQuery(req.KeyWords, "name", "goods_bried"))
	}
	if req.IsHot {
		q.Filter(elastic.NewTermQuery("is_hot", true))
	}
	if req.IsNew {
		q.Filter(elastic.NewTermQuery("is_new", true))
	}
	if req.PriceMin > 0 {
		q.Filter(elastic.NewRangeQuery("shop_price").Gte(req.PriceMin))
	}
	if req.PriceMax > 0 {
		q.Filter(elastic.NewRangeQuery("shop_price").Lte(req.PriceMax))
	}
	if req.Brand > 0 {
		q.Filter(elastic.NewTermQuery("brands_id", req.Brand))
	}
	var subQuery string
	categoryIds := make([]interface{}, 0)
	if req.TopCategory > 0 {
		var category model.Category
		if result := global.DB.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}
		if category.Level == 1 {
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category WHERE parent_category_id=%d) ", req.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("select id from category where parent_category_id=%d", req.TopCategory)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("select id from category where id=%d", req.TopCategory)
		}
		type Result struct {
			ID int32 `json:"id"`
		}
		var result []Result
		global.DB.Model(&model.Category{}).Raw(subQuery).Scan(&result)
		for _, re := range result {
			categoryIds = append(categoryIds, re.ID)
		}

		q = q.Filter(elastic.NewTermsQuery("category_id", categoryIds...))
	}

	if req.Pages == 0 {
		req.Pages = 1
	}
	switch {
	case req.PagePerNums > 100:
		req.PagePerNums = 100
	case req.PagePerNums <= 0:
		req.PagePerNums = 10
	}

	result, err := global.EsClient.Search().
		Index(model.EsGoods{}.GetIndexName()).
		Query(q).
		From(int(req.Pages)).
		Size(int(req.PagePerNums)).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	goodsListResponse.Total = int32(result.Hits.TotalHits.Value)

	goodsIds := make([]int32, 0)
	for _, value := range result.Hits.Hits {
		goods := model.EsGoods{}
		_ = json.Unmarshal(value.Source, &goods)
		goodsIds = append(goodsIds, goods.ID)

	}

	var goods []model.Goods
	r := localDB.Preload("Category").Preload("Brands").Find(&goods, goodsIds)
	if r.Error != nil {
		return nil, r.Error
	}

	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}

	return goodsListResponse, nil

	// 使用 mysql 搜索
	// var goods []model.Goods
	// goodsListResponse := &proto.GoodsListResponse{}
	// localDB := global.DB.Model(model.Goods{})
	// if req.KeyWords != "" {
	// 	localDB = localDB.Where("name LIKE ?", "%"+req.KeyWords+"%")
	// }
	// if req.IsHot {
	// 	localDB = localDB.Where(model.Goods{IsHot: true})
	// }
	// if req.IsNew {
	// 	localDB = localDB.Where(model.Goods{IsNew: true})
	// }
	// if req.PriceMin > 0 {
	// 	localDB = localDB.Where("shop_price>=?", req.PriceMin)
	// }
	// if req.PriceMax > 0 {
	// 	localDB = localDB.Where("shop_price<=?", req.PriceMax)
	// }
	// if req.Brand > 0 {
	// 	localDB = localDB.Where("brands_id=?", req.Brand)
	// }
	// var subQuery string
	// if req.TopCategory > 0 {
	// 	var category model.Category
	// 	if result := global.DB.First(&category, req.TopCategory); result.RowsAffected == 0 {
	// 		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	// 	}
	// 	if category.Level == 1 {
	// 		subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category WHERE parent_category_id=%d) ", req.TopCategory)
	// 	} else if category.Level == 2 {
	// 		subQuery = fmt.Sprintf("select id from category where parent_category_id=%d", req.TopCategory)
	// 	} else if category.Level == 3 {
	// 		subQuery = fmt.Sprintf("select id from category where id=%d", req.TopCategory)
	// 	}
	// 	localDB = localDB.Where(fmt.Sprintf("category_id in (%s)", subQuery))
	// }

	// var count int64
	// localDB.Count(&count)
	// goodsListResponse.Total = int32(count)

	// result := localDB.Preload("Category").Preload("Brands").Scopes(Paginate((int(req.Pages)), int(req.PagePerNums))).Find(&goods)
	// if result.Error != nil {
	// 	return nil, result.Error
	// }

	// for _, good := range goods {
	// 	goodsInfoResponse := ModelToResponse(good)
	// 	goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	// }

	// return goodsListResponse, nil
}

// 现在用户提交订单有多个商品，你得批量查询商品的信息吧
func (s *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	goodsListResponse := &proto.GoodsListResponse{}
	var goods []model.Goods
	result := global.DB.Where(req.Id).Find(&goods)
	for _, good := range goods {
		goodsInfoRespnse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoRespnse)
	}
	goodsListResponse.Total = int32(result.RowsAffected)
	return goodsListResponse, nil
}

func (s *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}
	//这里没有看到图片文件是如何上传， 在微服务中 普通的文件上传已经不再使用
	goods := model.Goods{
		Brands:          brand,
		BrandsId:        brand.ID,
		Category:        category,
		CategoryId:      category.ID,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.GoodsBrief,
		ShipFree:        req.ShipFree,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		OnSale:          req.OnSale,
	}

	tx := global.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &proto.GoodsInfoResponse{
		Id: goods.ID,
	}, nil
}

func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*empty.Empty, error) {
	if result := global.DB.Delete(&model.Goods{BaseModel: model.BaseModel{ID: req.Id}}, req.Id); result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	return &empty.Empty{}, nil
}

func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*empty.Empty, error) {
	var goods model.Goods

	if result := global.DB.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品不存在")
	}

	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	goods.Brands = brand
	goods.BrandsId = brand.ID
	goods.Category = category
	goods.CategoryId = category.ID
	goods.Name = req.Name
	goods.GoodsSn = req.GoodsSn
	goods.MarketPrice = req.MarketPrice
	goods.ShopPrice = req.ShopPrice
	goods.GoodsBrief = req.GoodsBrief
	goods.ShipFree = req.ShipFree
	goods.Images = req.Images
	goods.DescImages = req.DescImages
	goods.GoodsFrontImage = req.GoodsFrontImage
	goods.IsNew = req.IsNew
	goods.IsHot = req.IsHot
	goods.OnSale = req.OnSale

	tx := global.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &empty.Empty{}, nil
}

func (s *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var goods model.Goods

	if result := global.DB.Preload("Category").Preload("Brands").First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	goodsInfoResponse := ModelToResponse(goods)
	return &goodsInfoResponse, nil
}
