package model

import "gorm.io/gorm"

// OrderItems 订单商品项
type OrderItems struct {
	gorm.Model
	OrderId      int64   `gorm:"type:bigint;comment:订单ID"`
	ProductId    int64   `gorm:"type:bigint;comment:商品ID"`
	ProductName  string  `gorm:"type:varchar(200);comment:商品名称"`
	ProductImage string  `gorm:"type:varchar(500);comment:商品图片"`
	Price        float64 `gorm:"type:decimal(10,2);comment:商品单价"`
	Quantity     int     `gorm:"type:int;comment:购买数量"`
	Total        float64 `gorm:"type:decimal(10,2);comment:商品总金额"`
	//Status       int     `gorm:"type:tinyint;default:1;comment:状态: 1-正常, 2-已退款"`
	//ProductSku    string  `gorm:"type:varchar(100);comment:商品规格"`
	//OriginalPrice float64 `gorm:"type:decimal(10,2);comment:商品原价"`
	//PointsEarned int     `gorm:"type:int;default:0;comment:获得积分"`
}

func (i *OrderItems) FindOrderItemByOrderId(db *gorm.DB, id uint) error {
	return db.Debug().Where("order_id = ?", id).First(&i).Error
}
