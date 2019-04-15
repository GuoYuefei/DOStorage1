package main

import (
	"net/http"
	"storage/distributed/date/heartbeat"
	"storage/distributed/date/locate"
	"storage/distributed/date/objects"
	"storage/distributed/doslog"
)

func main() {
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	doslog.Loger.Printf("")
}
