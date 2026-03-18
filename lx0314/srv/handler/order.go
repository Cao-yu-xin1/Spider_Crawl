package handler

import (
	"context"
	"errors"
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/nacos"
	__ "github.com/Cao-yu-xin1/Spider_Crawl/lx0314/proto"
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/srv/model"
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/srv/pkg/alipay"
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/srv/pkg/ordersn"
	"gorm.io/gorm"
)

func (s *Server) CreateOrder(_ context.Context, in *__.CreateOrderReq) (*__.CreateOrderResp, error) {
	orderNo := ordersn.OrderNo()
	total := 0.0
	var orderItems []model.OrderItems
	for _, item := range in.OrderItems {
		var product model.Product
		err := product.FindProductById(nacos.DB, item.ProductId)
		if err != nil {
			return nil, errors.New("商品不存在")
		}
		if product.ID == 0 {
			return nil, errors.New("商品不存在")
		}
		total += product.Price * float64(item.Quantity)
		orderItems = append(orderItems, model.OrderItems{
			ProductId:    int64(product.ID),
			ProductName:  product.Name,
			ProductImage: product.Images,
			Price:        product.Price,
			Quantity:     int(item.Quantity),
			Total:        total,
		})
	}
	order := model.Order{
		Model:   gorm.Model{},
		OrderNo: orderNo,
		UserID:  in.UserId,
		Total:   total,
	}
	err := order.CreateOrder(nacos.DB)
	if err != nil {
		return nil, errors.New("创建订单失败")
	}
	for i, _ := range orderItems {
		orderItems[i].OrderId = int64(order.ID)
	}
	err = order.CreateOrderItem(nacos.DB, orderItems)
	if err != nil {
		return nil, errors.New("创建订单商品失败")
	}
	PayURL := alipay.AliPay(orderNo, total)
	return &__.CreateOrderResp{
		OrderNo: orderNo,
		Total:   float32(total),
		PayType: PayURL,
	}, nil
}

func (s *Server) NotifyPay(_ context.Context, in *__.NotifyPayReq) (*__.NotifyPayResp, error) {
	var order model.Order
	err := order.FindOrderByOrderSn(nacos.DB, in.OrderNo)
	if err != nil {
		return nil, errors.New("订单号不存在")
	}

	order.Status = 1
	//err = order.SaveOrder(nacos.DB)
	err = nacos.DB.Debug().Save(&order).Error
	if err != nil {
		return nil, errors.New("订单状态修改失败")
	}
	var orderItem []model.OrderItems
	//err = orderItem.FindOrderItemByOrderId(nacos.DB, order.ID)
	err = nacos.DB.Debug().Where("order_id = ?", order.ID).First(&orderItem).Error
	if err != nil {
		return nil, errors.New("订单详情不存在")
	}
	for _, items := range orderItem {
		var product model.Product
		err := product.FindProductById(nacos.DB, items.ProductId)
		if err != nil {
			return nil, errors.New("商品不存在")
		}
		//RabbitMQ.SendStockDeductMsg(strconv.Itoa(int(product.ID)), items.Quantity)
		product.Stock -= items.Quantity
		err = product.SaveProduct(nacos.DB)
		if err != nil {
			return nil, errors.New("库存扣减失败")
		}
	}
	return &__.NotifyPayResp{
		//OrderId: int64(order.ID),
	}, err
}
