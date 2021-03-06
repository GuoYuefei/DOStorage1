package main

import (
	"flag"
	"github.com/GuoYuefei/DOStorage1/distributed/config"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/heartbeat"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/locate"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/objects"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/temp"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/versions"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"net/http"
)

// 接口服务程序主函数
func main() {
	utils.Log.SetPriority(utils.Debug)
	configs()				// 配置文件
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	utils.Log.Println(utils.Info, "interface server will run in ", config.ServerInf.LISTEN_ADDRESS)
	utils.FailOnError(http.ListenAndServe(config.ServerInf.LISTEN_ADDRESS, nil), "Fail to open a server")
}

func configs() {
	config.Flags(config.TypeSInf)
	flag.Parse()
	config.ConfigParse(config.TypeSInf)
}
