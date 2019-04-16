package main

import (
	"net/http"
	"os"
	"storage/distributed/data/heartbeat"
	"storage/distributed/data/locate"
	"storage/distributed/data/objects"
	"storage/distributed/doslog"
)

func main() {
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	doslog.FailOnError(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil),
		"Fail to open a data server")
}
