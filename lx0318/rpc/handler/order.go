package handler

import (
	"context"
	"errors"
	"lx0318/config"
	"lx0318/pkg/alipay"
	"lx0318/pkg/ordersn"
	__ "lx0318/proto"
	"lx0318/rpc/model"
)

func (s *Server) CreateOrder(_ context.Context, in *__.CreateOrderReq) (*__.CreateOrderResp, error) {
	orderNo := ordersn.OrderSn()
	total := 0.0
	var orderItems []model.OrderItem
	for _, item := range in.Item {
		var product model.Product
		err := product.FindProductById(config.DB, item.ProductId)
		if err != nil {
			return nil, errors.New("商品不存在")
		}
		if product.ID == 0 {
			return nil, errors.New("商品不存在")
		}
		total += product.Price * float64(item.Quantity)
		orderItems = append(orderItems, model.OrderItem{
			ProductId:    int(product.ID),
			ProductName:  product.Name,
			ProductPrice: product.Price,
			Quantity:     int(item.Quantity),
		})
	}
	order := model.Order{
		OrderNo:   orderNo,
		MemberId:  int(in.MemberId),
		AddressId: int(in.AddressId),
		Total:     total,
	}
	err := order.CreateOrder(config.DB)
	if err != nil {
		return nil, errors.New("订单创建失败")
	}
	for i, _ := range orderItems {
		orderItems[i].OrderId = int(order.ID)
	}
	err = order.CreateOrderItem(config.DB, orderItems)
	if err != nil {
		return nil, errors.New("订单详情创建失败")
	}
	aliPay := alipay.AliPay(orderNo, total)
	return &__.CreateOrderResp{
		OrderNo: orderNo,
		Total:   float32(total),
		PayUrl:  aliPay,
	}, nil
}

func (s *Server) NotifyPay(_ context.Context, in *__.NotifyPayReq) (*__.NotifyPayResp, error) {
	var order model.Order
	err := order.FindOrderByOrderNo(config.DB, in.OrderNo)
	if err != nil {
		return nil, errors.New("订单号不存在")
	}
	order.Status = 1
	err = order.SaveOrder(config.DB)
	if err != nil {
		return nil, errors.New("修改订单状态失败")
	}
	var orderItems []model.OrderItem
	err = config.DB.Debug().Where("order_id = ?", order.ID).Find(&orderItems).Error
	if err != nil {
		return nil, errors.New("订单详情不存在")
	}
	for _, item := range orderItems {
		var product model.Product
		err := product.FindProductById(config.DB, int64(item.ProductId))
		if err != nil {
			return nil, errors.New("商品不存在")
		}
		if product.ID == 0 {
			return nil, errors.New("商品不存在")
		}
		product.Stock -= item.Quantity
		err = product.SaveProduct(config.DB)
		if err != nil {
			return nil, errors.New("库存扣减失败")
		}
	}

	return &__.NotifyPayResp{}, nil
}
