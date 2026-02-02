package catalogue

import "os"

var str = []string{
	"bff-gateway",
	"bff-gateway/basic",
	"bff-gateway/basic/cmd",
	"bff-gateway/basic/init",
	"bff-gateway/basic/config",
	"bff-gateway/basic/proto",
	"bff-gateway/handler",
	"bff-gateway/model",
	"bff-gateway/router",
	"bff-gateway/pkg",
	"rpc-sev",
	"rpc-sev/basic",
	"rpc-sev/basic/cmd",
	"rpc-sev/basic/init",
	"rpc-sev/basic/config",
	"rpc-sev/basic/proto",
	"rpc-sev/handler",
	"rpc-sev/model",
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
