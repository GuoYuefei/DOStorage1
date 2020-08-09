package main

import (
	"github.com/GuoYuefei/DOStorage1/distributed/doslog"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/heartbeat"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/locate"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/objects"
	"net/http"
	"os"
)

// 接口服务程序主函数
func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Put)
	http.HandleFunc("/locate/", locate.Handler)
	doslog.FailOnError(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil), "Fail to open a server")
}
