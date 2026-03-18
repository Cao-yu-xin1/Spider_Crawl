package init

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"lx0318/config"
	"lx0318/rpc/model"
	"sync"
	"time"
)

func MysqlInit() {
	cng := config.GlobalConfig.Mysql
	once := sync.Once{}
	var err error
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cng.User, cng.Password, cng.Host, cng.Port, cng.Database)
	once.Do(func() {
		config.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("mysql 连接失败" + err.Error())
		}
		fmt.Println("mysql 连接成功")
	})
	err = config.DB.AutoMigrate(
		&model.MemberRegister{},
		&model.MemberInfo{},
		&model.Points{},
		&model.Product{},
		&model.Stock{},
		&model.Logistics{},
		&model.OrderItem{},
		&model.Order{},
	)
	if err != nil {
		panic("mysql 自动迁移失败" + err.Error())
	}
	fmt.Println("mysql 自动迁移成功")
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := config.DB.DB()

	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
}
