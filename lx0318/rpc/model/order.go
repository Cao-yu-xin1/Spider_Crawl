package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	OrderNo   string  `gorm:"type:varchar(30);comment:订单号"`
	MemberId  int     `gorm:"type:int;comment:会员id"`
	AddressId int     `gorm:"type:int;comment:地址id"`
	Total     float64 `gorm:"type:decimal(10,2);comment:总价"`
	Status    int     `gorm:"type:int;default:0;comment:状态:0-待付款 1-已付款 2-已发货 3-已完成 4-已取消"`
	//PayUrl    string  `gorm:"type:varchar(30);comment:订单号"`
}

func (o *Order) CreateOrder(db *gorm.DB) error {
	return db.Debug().Create(&o).Error
}

func (o *Order) CreateOrderItem(db *gorm.DB, items []OrderItem) error {
	return db.Debug().Create(&items).Error
}

func (o *Order) FindOrderByOrderNo(db *gorm.DB, orderNo string) error {
	return db.Debug().Where("order_no = ?", orderNo).First(&o).Error
}

func (o *Order) SaveOrder(db *gorm.DB) error {
	return db.Debug(). /*.Model(&Order{})*/ Save(&o).Error
}

type OrderItem struct {
	gorm.Model
	OrderId      int     `gorm:"type:int;comment:订单id"`
	ProductId    int     `gorm:"type:int;comment:商品id"`
	ProductName  string  `gorm:"type:varchar(30);comment:商品名称"`
	ProductPrice float64 `gorm:"type:decimal(10,2);comment:商品价格"`
	Quantity     int     `gorm:"type:int;comment:商品名称"`
}
