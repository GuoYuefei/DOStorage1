package main

import (
	"./objects"
	"log"
	"net/http"
	"os"
)

func main() {
	var listenAddr string = "LISTEN_ADDRESS"
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv(listenAddr), nil))
}
