package config

import (
	"gorm.io/gorm"
	__ "lx0318/proto"
)

var (
	GlobalConfig  AppConfig
	DB            *gorm.DB
	ServiceClient __.ServiceClient
)
