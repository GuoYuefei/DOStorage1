package main

import (
	"net/http"
	"os"
	"storage/distributed/doslog"
	"storage/distributed/interface/heartbeat"
	"storage/distributed/interface/locate"
	"storage/distributed/interface/objects"
)

// 接口服务程序主函数
func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Put)
	http.HandleFunc("/locate/", locate.Handler)
	doslog.FailOnError(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil), "Fail to open a server")
}
