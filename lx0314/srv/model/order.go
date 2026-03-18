package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	OrderNo   string  `gorm:"type:varchar(30);comment:订单号"`
	UserID    int64   `gorm:"type:bigint(20);comment:用户ID"`
	Total     float64 `gorm:"type:decimal(10,2);comment:订单总金额"`
	PayAmount float64 `gorm:"type:decimal(10,2);comment:订单支付金额"`
	Status    int8    `gorm:"type:tinyint(4);default:0;comment:订单状态,0-待支付,1-已支付,2-已取消"`
	PayType   int8    `gorm:"type:tinyint(4);default:0;comment:支付方式,0-支付宝,1-微信"`
	PayUrl    int8    `gorm:"type:tinyint(4);default:0;comment:支付URL"`
}

func (o *Order) CreateOrder(db *gorm.DB) error {
	return db.Debug().Create(&o).Error
}

func (o *Order) CreateOrderItem(db *gorm.DB, orderItems []OrderItems) error {
	return db.Debug().Create(&orderItems).Error
}

func (o *Order) FindOrderByOrderSn(db *gorm.DB, no string) error {
	return db.Debug().Where("order_no = ?", no).Limit(1).Find(&o).Error
}

func (o *Order) SaveOrder(db *gorm.DB) error {
	return db.Debug().Model(&Order{}).Save(&o).Error
}
