package model

import "gorm.io/gorm"

// User 用户表
type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);comment:用户名"`
	Email    string `gorm:"type:varchar(100);comment:邮箱"`
	Password string `gorm:"type:varchar(255);comment:密码"`
	Phone    string `gorm:"type:varchar(20);comment:手机号"`
	Avatar   string `gorm:"type:varchar(255);comment:头像"`
	Status   int8   `gorm:"type:tinyint;default:1;comment:状态"`
	Level    int8   `gorm:"type:tinyint;default:0;comment:会员等级"`
}

func (u *User) FindUserByUsername(db *gorm.DB, username string) error {
	return db.Debug().Where("username = ?", username).Limit(1).Find(&u).Error
}

func (u *User) CreateUser(db *gorm.DB) error {
	return db.Debug().Create(&u).Error
}
