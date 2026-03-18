package main

import "os"

var str = []string{
	"bff-api",
	"bff-api/basic",
	"bff-api/basic/cmd",
	"bff-api/basic/init",
	"bff-api/basic/config",
	"bff-api/basic/proto",
	"bff-api/handler",
	"bff-api/handler/request",
	"bff-api/handler/response",
	"bff-api/handler/service",
	"bff-api/middleware",
	"bff-api/router",
	"bff-api/pkg",
	"srv",
	"srv/basic",
	"srv/basic/cmd",
	"srv/basic/init",
	"srv/basic/config",
	"srv/basic/proto",
	"srv/pkg",
	"srv/handler",
	"srv/model",
}

func CreateDir() {
	for _, s := range str {
		err := os.MkdirAll(s, os.ModePerm)
		if err != nil {
			panic("目录创建失败")
		}
		println("目录创建成功")
	}
}

func main() {
	CreateDir()
}
