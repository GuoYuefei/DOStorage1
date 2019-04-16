package objects

import (
	"io"
	"log"
	"net/http"
	"os"
	"storage/distributed/doslog"
	"strings"
)

var storage_root string = "STORAGE_ROOT"

func put(w http.ResponseWriter, r *http.Request) {
	f, e := os.Create(os.Getenv(storage_root)+"/objects/"+strings.Split(r.URL.EscapedPath(),"/")[2])

	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(f, r.Body)
}

func get(w http.ResponseWriter, r *http.Request) {
	f, e := os.Open(os.Getenv(storage_root)+"/objects/"+strings.Split(r.URL.EscapedPath(), "/")[2])

	if e != nil {
		doslog.FailOnError(e, "not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	io.Copy(w, f)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPut {
		put(w, r)
		return
	}

	if m == http.MethodGet {
		get(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}