package main

import (
	"github.com/GuoYuefei/DOStorage1/distributed/data/heartbeat"
	"github.com/GuoYuefei/DOStorage1/distributed/data/locate"
	"github.com/GuoYuefei/DOStorage1/distributed/data/objects"
	"github.com/GuoYuefei/DOStorage1/distributed/doslog"
	"net/http"
	"os"
)

func main() {
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	doslog.FailOnError(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil),
		"Fail to open a data server")
}
