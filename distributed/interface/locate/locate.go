// 主要用于处理get请求
package locate

import (
	"encoding/json"
	"github.com/GuoYuefei/DOStorage1/distributed/config"
	"github.com/GuoYuefei/DOStorage1/distributed/es"
	"github.com/GuoYuefei/DOStorage1/distributed/rabbitmq"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// copy from object get
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	versionId := r.URL.Query()["version"]
	version := 0
	var e error
	if len(versionId) != 0 {
		version, e = strconv.Atoi(versionId[0])
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	utils.Log.Println(utils.Debug, "get for ", name, "version is ", version)
	meta, e := es.GetMetadata(name, version)
	utils.Log.Println(utils.Debug, "get metadata is ", meta)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if meta.Hash == "" {
		utils.Log.Println(utils.Info, "meta's hash is none")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	object := url.PathEscape(meta.Hash)
	// end copy from object get
	info := Locate(object)
	if len(info) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, _ := json.Marshal(info)
	_, _ = w.Write(b)
}

func Locate(hash string) string {
	q := rabbitmq.New(config.Pub.RABBITMQ_SERVER)
	q.Publish("dataServers", hash)
	c := q.Consume()
	go func() {
		time.Sleep(1*time.Second)
		q.Close()
	}()
	msg := <-c
	s, _ := strconv.Unquote(string(msg.Body))
	return s
}

func Exist(name string) bool {
	return Locate(name) != ""
}
