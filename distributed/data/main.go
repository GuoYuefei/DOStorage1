package main

import (
	"flag"
	"github.com/GuoYuefei/DOStorage1/distributed/config"
	"github.com/GuoYuefei/DOStorage1/distributed/data/heartbeat"
	"github.com/GuoYuefei/DOStorage1/distributed/data/locate"
	"github.com/GuoYuefei/DOStorage1/distributed/data/objects"
	"github.com/GuoYuefei/DOStorage1/distributed/data/temp"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"net/http"
)

func main() {
	utils.Log.SetPriority(utils.Debug)
	configs()
	utils.Log.Println(utils.Info, "data server STORAGE_ROOT is ", config.ServerData.STORAGE_ROOT)
	locate.CollectObjects() 			// 启动时先收集本地的对象，存入内存
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	utils.Log.Println(utils.Info, "data server will run in ", config.ServerData.LISTEN_ADDRESS)
	utils.PanicOnError(http.ListenAndServe(config.ServerData.LISTEN_ADDRESS, nil),
		"Fail to open a data server")
}

func configs() {
	config.Flags(config.TypeSData)
	flag.Parse()
	config.ConfigParse(config.TypeSData)
}
