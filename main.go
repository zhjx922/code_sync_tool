package main

import (
	"code_sync_tool/config"
	"code_sync_tool/event"
	"code_sync_tool/helper"
	log "code_sync_tool/log"
	"code_sync_tool/watcher"
	"flag"
)

// 初始化
func init() {
	log.Println("初始化中……")
	helper.SetLimit()
}

func main() {
	// 配置文件处理
	var source string
	flag.StringVar(&source, "c", "./conf.ini", "配置文件路径")
	flag.Parse()

	e := make(chan event.FileEvent)
	wa := watcher.New(e)
	configs := config.GetConfigsByIni(source)

	// 添加监控目录
	for _, c := range configs {
		wa.AddDir(c.LocalPath)
	}
	log.Println("开始……")

	event.Run(configs, e)

	do := make(chan bool)

	<-do
}
