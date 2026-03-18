package model

import "gorm.io/gorm"

// Logistics 物流信息
type Logistics struct {
	gorm.Model
	OrderId        int64  `gorm:"type:bigint;unique;comment:订单ID"`
	OrderNo        string `gorm:"type:varchar(30);comment:订单号"`
	UserId         int64  `gorm:"type:bigint;comment:用户ID"`
	ExpressCompany string `gorm:"type:varchar(100);comment:快递公司"`
	ExpressNo      string `gorm:"type:varchar(100);comment:快递单号"`
	ReceiverName   string `gorm:"type:varchar(50);comment:收货人姓名"`
	ReceiverPhone  string `gorm:"type:varchar(20);comment:收货人电话"`
	Province       string `gorm:"type:varchar(50);comment:省"`
	City           string `gorm:"type:varchar(50);comment:市"`
	Area           string `gorm:"type:varchar(50);comment:区/县"`
	Address        string `gorm:"type:varchar(500);comment:详细地址"`
	Status         int    `gorm:"type:tinyint;default:0;comment:状态: 0-待发货, 1-已发货, 2-已签收, 3-已退回"`
	LogisticsId    int64  `gorm:"type:bigint;comment:物流ID"`
}
