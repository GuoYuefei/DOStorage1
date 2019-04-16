package objects

import (
	"fmt"
	"io"
	"net/http"
	"storage/distributed/doslog"
	"storage/distributed/interface/heartbeat"
	"storage/distributed/interface/objectstream"
	"strings"
)

func Put(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2] // todo 2
	c, e := storeObject(r.Body, object)
	doslog.FailOnError(e, "Fail to storeObject")
	w.WriteHeader(c)

}

func storeObject(r io.Reader, object string) (int, error) {
	stream, e := putStream(object)
	if e != nil {
		doslog.FailOnError(e, e.Error())
		return http.StatusServiceUnavailable, e
	}

	io.Copy(stream, r)
	e = stream.Close()
	if e != nil {
		doslog.FailOnError(e, "Fail ot close the stream")
		return http.StatusInternalServerError, e
	}
	return http.StatusOK, nil
}

func putStream(object string) (*objectstream.PutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}

	return objectstream.NewPutStream(server, object), nil
}
