package service

import (
	"fmt"
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/mq/sc"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NotifyPay(c *gin.Context) {
	c.Request.ParseForm()
	fmt.Println("这是NotifyPay", c.Request.PostForm)
	/*
		这是NotifyPay map[app_id:[9021000157679181]
		auth_app_id:[9021000157679181]
		buyer_id:[2088722087798772]
		buyer_pay_amount:[5.50]
		charset:[utf-8]
		fund_bill_list:[[{"amount":"5.50","fundChannel":"ALIPAYACCOUNT"}]]
		gmt_create:[2026-03-16 16:59:43]
		gmt_payment:[2026-03-16 16:59:53]
		invoice_amount:[5.50]
		notify_id:[2026031601222165954098770507586668]
		notify_time:[2026-03-16 16:59:55]
		notify_type:[trade_status_sync]
		out_trade_no:[202603161659239154d0b7]
		point_amount:[0.00]
		receipt_amount:[5.50]
		seller_id:[2088721087767621]
		sign:[cLKnU4GKDbh1CUabAfbiWsGKERmjeHUfZWsVxvur76HnBR7csfaIxLfrTvBh9VMSDNg8jF42jrp8J5PUIAr3gEIyJvk2AE152WQUYFaiepfIEVLZOnWd7l4swXO1SL1Ngv0n7H16XYama1Re6s23zgXmVOkm0tzv6MbDS8+OmvuP6c1F9fVFZ20uUgrzugIKrSXUVuYNlRRcgk1Ogwa0IUwTuppmS+LcbdD7lGECtEzsthQ2KQZA6XvfL0vie1mivls43z0euu2R8IRS9KuGyc6/04qqrjskVeM4d6nUnWyTHln4vV8tjzGcRdB0vW8Zs2I9dVk8av0x4AoI8gDpfQ==]
		sign_type:[RSA2]
		subject:[标题]
		total_amount:[5.50]
		trade_no:[2026031622001498770507762375]
		trade_status:[TRADE_SUCCESS]
		version:[1.0]]
	*/

	/*
		out_trade_no:[202603161659239154d0b7]
		trade_status:[TRADE_SUCCESS]
	*/

	tradeStatus := c.Request.PostForm.Get("trade_status")

	if tradeStatus != "TRADE_SUCCESS" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "交易失败",
		})
		return
	}

	outTradeNo := c.Request.PostForm.Get("out_trade_no")

	if outTradeNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "订单号不存在",
		})
		return
	}

	go func() {
		sc.SendMessage("topic", outTradeNo)
	}()

	//_, err := config.ServiceClient.NotifyPay(c, &__.NotifyPayReq{OrderNo: outTradeNo})
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"code": 400,
	//		"msg":  "订单修改失败",
	//	})
	//	return
	//}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		//"data": gin.H{
		//	"notifyPay": notifyPay,
		//},
	})
	return

	//
	//var order model.Order
	//err := nacos.DB.Debug().Where("order_sn = ?", outTradeNo).First(&order).Error
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"code": 400,
	//		"msg":  "订单号不存在",
	//	})
	//	return
	//}
	//
	//order.Status = 2
	//
	//err = nacos.DB.Debug().Save(&order).Error
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"code": 400,
	//		"msg":  "订单状态修改失败",
	//	})
	//	return
	//}
	//
	//var orderItem []model.OrderItems
	//err = nacos.DB.Debug().Where("order_id = ?", order.ID).First(&orderItem).Error
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"code": 400,
	//		"msg":  "订单详情不存在",
	//	})
	//}
	//
	//for _, item := range orderItem {
	//	var product model.Product
	//
	//	err = nacos.DB.Debug().Where("id = ?", item.ProductId).First(&product).Error
	//	if err != nil {
	//		c.JSON(http.StatusBadRequest, gin.H{
	//			"code": 400,
	//			"msg":  "商品不存在",
	//		})
	//		return
	//	}
	//
	//	product.Stock -= item.Quantity
	//
	//	err = nacos.DB.Debug().Save(&product).Error
	//	if err != nil {
	//		c.JSON(http.StatusBadRequest, gin.H{
	//			"code": 400,
	//			"msg":  "库存扣减失败",
	//		})
	//		return
	//	}
	//}
}
