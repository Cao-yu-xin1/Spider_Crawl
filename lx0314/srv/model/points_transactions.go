package model

import "gorm.io/gorm"

// PointsTransactions 积分变动记录
type PointsTransactions struct {
	gorm.Model
	UserId      int64  `gorm:"type:bigint;comment:用户ID"`
	Points      int64  `gorm:"type:bigint;comment:积分变动数量(正数增加,负数减少)"`
	Type        int8   `gorm:"type:tinyint;comment:类型: 1-获得, 2-消费, 3-过期, 4-管理员调整"`
	RelatedType int8   `gorm:"type:tinyint;comment:关联类型: 1-订单, 2-评价, 3-签到, 4-活动"`
	RelatedId   int64  `gorm:"type:bigint;comment:关联ID(订单ID/评价ID等)"`
	Balance     int64  `gorm:"type:bigint;comment:变动后总积分"`
	Description string `gorm:"type:varchar(200);comment:描述"`
}
