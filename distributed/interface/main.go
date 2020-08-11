package main

import (
	"github.com/GuoYuefei/DOStorage1/distributed/config"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/heartbeat"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/locate"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/objects"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"net/http"
)

// 接口服务程序主函数
func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	utils.FailOnError(http.ListenAndServe(config.ServerInf.LISTEN_ADDRESS, nil), "Fail to open a server")
}
