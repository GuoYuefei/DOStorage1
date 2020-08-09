package main

import (
	"github.com/GuoYuefei/DOStorage1/singlemachine/objects"
	"log"
	"net/http"
	"os"
)

func main() {
	var listenAddr string = "LISTEN_ADDRESS"
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv(listenAddr), nil))
}
