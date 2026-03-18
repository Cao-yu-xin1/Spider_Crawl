package model

import "gorm.io/gorm"

// Inventory 库存信息
type Inventory struct {
	gorm.Model
	ProductId      int64  `gorm:"type:bigint;comment:商品ID"`
	TotalStock     int    `gorm:"type:int;default:0;comment:总库存"`
	FrozenStock    int    `gorm:"type:int;default:0;comment:冻结库存(下单占用)"`
	AvailableStock int    `gorm:"type:int;default:0;comment:可用库存"`
	WarehouseId    int64  `gorm:"type:bigint;comment:仓库ID"`
	Location       string `gorm:"type:varchar(100);comment:库位"`
}

//git config --global user.name "caoyuxin"
//git config --global user.email "1041862764@qq.com"
