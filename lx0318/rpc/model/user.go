package model

import (
	"gorm.io/gorm"
	__ "lx0318/proto"
)

// MemberRegister 注册会员表
type MemberRegister struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);comment:用户名"`
	Password string `gorm:"type:varchar(255);comment:密码"` //json:"-" 表示不返回密码
	Phone    string `gorm:"type:varchar(20);comment:手机号"`
}

// MemberInfo 会员信息表
type MemberInfo struct {
	gorm.Model
	RegisterID uint64 `gorm:"type:int;comment:注册id"`
	Nickname   string `gorm:"type:varchar(50);comment:昵称"`
	Level      int    `gorm:"type:tinyint;default:1;comment:等级"`
}

// Points 积分明细表
type Points struct {
	gorm.Model
	MemberID uint64 `gorm:"type:int;comment:会员id"`
	Points   int    `gorm:"type:int;comment:积分数值"`
	Type     int    `gorm:"type:tinyint;comment:类型：1-收入 2-支出"` // 1-收入 2-支出
}

// Product 商品表
type Product struct {
	gorm.Model
	Name   string  `gorm:"type:varchar(200);comment:商品名称"`
	Price  float64 `gorm:"type:decimal(10,2);comment:价格"`
	Stock  int     `gorm:"type:int;default:0;comment:库存"`               // 0-下架 1-上架
	Status int     `gorm:"type:tinyint;default:1;comment:状态：0-下架 1-上架"` // 0-下架 1-上架
}

func (p *Product) CreateProduct(db *gorm.DB) error {
	return db.Debug().Create(&p).Error
}

func (p *Product) FindProductById(db *gorm.DB, id int64) error {
	return db.Debug().Where("id = ?", id).Limit(1).Find(&p).Error
}

func (p *Product) DeleteProduct(db *gorm.DB) error {
	return db.Debug().Delete(&p).Error
}

func (p *Product) UpdateProduct(db *gorm.DB, id int64) error {
	return db.Debug().Where("id = ?", id).Updates(&p).Error
}

func (p *Product) SearchProduct(db *gorm.DB, in *__.SearchProductsReq) (list []*__.Products, err error) {
	if in.Page == 0 {
		in.Page = 1
	}
	if in.Size == 0 || in.Size > 100 {
		in.Size = 10
	}
	query := db.Debug().Model(&Product{})
	if in.Name != "" {
		query = query.Where("name LIKE ?", "%"+in.Name+"%")
	}
	if in.Status != 0 {
		query = query.Where("status = ?", in.Status)
	}
	offset := (in.Page - 1) * in.Size
	// 执行查询
	err = query.Offset(int(offset)).Limit(int(in.Size)).Find(&list).Error
	return list, nil
}

func (p *Product) SaveProduct(db *gorm.DB) error {
	return db.Debug().Save(&p).Error
}

// Stock 库存表
type Stock struct {
	gorm.Model
	ProductID uint64 `gorm:"type:int;comment:商品ID"`
	Quantity  int    `gorm:"type:int;default:0;comment:库存数量"`
}

// Logistics 物流配送表
type Logistics struct {
	gorm.Model
	OrderNo string `gorm:"type:varchar(50);comment:订单号"`
	//MemberId      int    `gorm:"type:int;comment:会员id"`
	LogisticsNo   string `gorm:"type:varchar(50);comment:物流单号"`
	Carrier       string `gorm:"type:varchar(50);comment:快递公司"`
	ReceiverName  string `gorm:"type:varchar(50);comment:收货人"`
	ReceiverPhone string `gorm:"type:varchar(20);comment:收货电话"`
	Status        int    `gorm:"type:tinyint;default:0;comment:状态：0-待发货 1-已发货 2-已签收"` // 0-待发货 1-已发货 2-已签收
}
