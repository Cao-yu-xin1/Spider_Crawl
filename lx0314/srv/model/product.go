package model

import (
	"gorm.io/gorm"
	__ "lx0314/proto"
)

// Product 商品表
type Product struct {
	gorm.Model
	Name     string  `gorm:"type:varchar(200);not null;comment:商品名称"`
	Subtitle string  `gorm:"type:varchar(200);comment:副标题"`
	Desc     string  `gorm:"type:text;comment:商品描述"`
	Images   string  `gorm:"type:text;comment:商品图片(JSON数组)"`
	Price    float64 `gorm:"type:decimal(10,2);not null;comment:当前价格"`
	Stock    int     `gorm:"type:int;default:0;comment:库存数量"`
	Status   int     `gorm:"type:tinyint;default:1;comment:状态: 0-下架, 1-上架, 2-删除"`
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
