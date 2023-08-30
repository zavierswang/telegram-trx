package main

import (
	"telegram-trx/pkg/bootstrap"
	"telegram-trx/pkg/core/cst"
)

func main() {
	// 初始配置文件
	bootstrap.LoadConfig(cst.AppName)
	// 初始化database
	bootstrap.ConnectDB()
	// 初始化出T私钥
	bootstrap.StoreKey()
	// 初始化Telegram
	err := bootstrap.Telegram()
	if err != nil {
		panic(err)
	}
}
