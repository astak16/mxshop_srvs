package handler

import (
	"context"
	"math/rand"

	"fmt"
	"mxshow_srvs/order_srv/global"
	"mxshow_srvs/order_srv/model"
	"mxshow_srvs/order_srv/proto"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}

func GenerateOrderSn(userId int32) string {
	// 生成规则 年月日时分秒 + 用户 id + 随机数
	now := time.Now()
	// rand.Seed(time.Now().UnixNano())
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(), userId, r.Intn(90)+10)
	return orderSn
}

func (*OrderServer) CartItemList(ctx context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	var shopCarts []model.ShoppingCart
	var rsp proto.CartItemListResponse
	if result := global.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&shopCarts); result.Error != nil {
		return nil, result.Error
	} else {
		rsp.Total = int32(result.RowsAffected)
	}
	for _, shopCart := range shopCarts {
		rsp.Data = append(rsp.Data, &proto.ShopCartInfoResponse{
			Id:      shopCart.ID,
			UserId:  shopCart.User,
			GoodsId: shopCart.Goods,
			Nums:    shopCart.Nums,
			Checked: shopCart.Checked,
		})
	}
	return &rsp, nil
}

func (*OrderServer) CreateCartItem(ctx context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	var shopCart model.ShoppingCart

	if result := global.DB.Where(&model.ShoppingCart{Goods: req.GoodsId, User: req.UserId}).First(&shopCart); result.RowsAffected == 1 {
		shopCart.Nums += req.Nums
	} else {
		shopCart.User = req.UserId
		shopCart.Goods = req.GoodsId
		shopCart.Nums = req.Nums
		shopCart.Checked = false
	}
	global.DB.Save(&shopCart)
	return &proto.ShopCartInfoResponse{
		Id: shopCart.ID,
	}, nil
}

func (*OrderServer) UpdateCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	var shopCart model.ShoppingCart
	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).First(&shopCart); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	shopCart.Checked = req.Checked
	if req.Nums > 0 {
		shopCart.Nums = req.Nums
	}
	global.DB.Save(&shopCart)
	return &emptypb.Empty{}, nil
}

func (*OrderServer) DeleteCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	return &emptypb.Empty{}, nil
}

func (*OrderServer) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	var orders []model.OrderInfo
	var rsp proto.OrderListResponse

	var total int64
	global.DB.Where(&model.OrderInfo{User: req.UserId}).Count(&total)
	rsp.Total = int32(total)

	global.DB.Scopes(Paginate(int(req.Page), int(req.PagePerNums))).Where(&model.OrderInfo{User: req.UserId}).Find(&orders)
	for _, order := range orders {
		rsp.Data = append(rsp.Data, &proto.OrderInfoResponse{
			Id:      order.ID,
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SignerMobile,
			AddTime: order.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &rsp, nil
}

func (*OrderServer) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	var order model.OrderInfo
	var rsp proto.OrderInfoDetailResponse
	if result := global.DB.Where(&model.OrderInfo{BaseModel: model.BaseModel{ID: req.Id}, User: req.UserId}).First(&order); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}

	orderInfo := proto.OrderInfoResponse{}
	orderInfo.Id = order.ID
	orderInfo.UserId = order.User
	orderInfo.OrderSn = order.OrderSn
	orderInfo.PayType = order.PayType
	orderInfo.Status = order.Status
	orderInfo.Post = order.Post
	orderInfo.Total = order.OrderMount
	orderInfo.Address = order.Address
	orderInfo.Name = order.SignerName
	orderInfo.Mobile = order.SignerMobile

	rsp.OrderInfo = &orderInfo

	var orderGoods []model.OrderGoods
	if result := global.DB.Where(&model.OrderGoods{Order: req.Id}).Find(&orderGoods); result.Error != nil {
		return nil, result.Error
	}

	for _, orderGood := range orderGoods {
		rsp.Goods = append(rsp.Goods, &proto.OrderItemResponse{
			GoodsId:    orderGood.Goods,
			GoodsName:  orderGood.GoodsName,
			GoodsPrice: orderGood.GoodsPrice,
			Nums:       orderGood.Nums,
		})
	}

	return &rsp, nil
}

func (*OrderServer) Create(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {

	// 1. 从购物车中获取到选中的商品
	// 2. 商品的价格自己查询 - 访问商品服务（跨微服务）
	// 3. 库存的扣减 - 访问库存服务（跨微服务）
	// 4. 订单的基本信息表 - 订单的商品信息表
	// 5. 从购物车中删除已购买的记录

	var goodsIds []int32
	var shopCarts []model.ShoppingCart
	goodsNumsMap := make(map[int32]int32)
	if result := global.DB.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Find(&shopCarts); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有选中结算的商品")
	}

	for _, shopCart := range shopCarts {
		goodsIds = append(goodsIds, shopCart.Goods)
		goodsNumsMap[shopCart.Goods] = shopCart.Nums
	}

	// 商品微服务
	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: goodsIds})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "批量查询商品信息失败")
	}

	var orderAmount float32
	var orderGoods []*model.OrderGoods
	var goodsInvInfo []*proto.GoodsInvInfo
	for _, goods := range goods.Data {
		orderAmount += goods.ShopPrice * float32(goodsNumsMap[goods.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      goods.Id,
			GoodsName:  goods.Name,
			GoodsImage: goods.GoodsFrontImage,
			GoodsPrice: goods.ShopPrice,
			Nums:       goodsNumsMap[goods.Id],
		})
		goodsInvInfo = append(goodsInvInfo, &proto.GoodsInvInfo{
			GoodsId: goods.Id,
			Nums:    goodsNumsMap[goods.Id],
		})
	}

	// 库存微服务
	if _, err := global.InventorySrvClient.Sell(context.Background(), &proto.SellInfo{GoodsInfo: goodsInvInfo}); err != nil {
		return nil, status.Errorf(codes.ResourceExhausted, "扣减库存失败")
	}

	// 生成订单表
	tx := global.DB.Begin()
	order := model.OrderInfo{
		OrderSn:      GenerateOrderSn(req.UserId),
		OrderMount:   orderAmount,
		Address:      req.Address,
		SignerName:   req.Name,
		SignerMobile: req.Mobile,
		Post:         req.Post,
		User:         req.UserId,
	}
	if result := tx.Save(&order); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, "创建订单失败")
	}

	for _, orderGood := range orderGoods {
		orderGood.Order = order.ID
	}
	if result := tx.CreateInBatches(&orderGoods, 100); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, "创建订单失败")
	}

	if result := tx.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, "创建订单失败")
	}
	tx.Commit()

	return &proto.OrderInfoResponse{
		Id:      order.ID,
		OrderSn: order.OrderSn,
		Total:   order.OrderMount,
	}, nil
}

func (*OrderServer) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*emptypb.Empty, error) {
	if result := global.DB.Model(&model.OrderInfo{}).Where("order_sn = ?", req.OrderSn).Update("status", req.Status); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "创建订单失败")
	}
	return &emptypb.Empty{}, nil
}
