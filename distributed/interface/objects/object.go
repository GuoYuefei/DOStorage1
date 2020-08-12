package objects

import (
	"fmt"
	"github.com/GuoYuefei/DOStorage1/distributed/es"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/heartbeat"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/locate"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/objectstream"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/util"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodGet {
		get(w, r)
		return
	}

	if m == http.MethodPut {
		put(w, r)
		return
	}

	if m == http.MethodDelete {
		del(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func del(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	version, e := es.SearchLatestVersion(name)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	e = es.PutMetadata(name, version.Version+1, 0, "")
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func put(w http.ResponseWriter, r *http.Request) {
	hash := util.GetHashFromHeader(r.Header)
	if hash == "" {
		utils.Log.Println(utils.Record, "missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	size := util.GetSizeFromHeader(r.Header)

	utils.Log.Println(utils.Debug, "from header, size is ", size)

	c, e := storeObject(r.Body, url.PathEscape(hash), size)
	if e != nil {
		utils.Log.Println(utils.Err, e)
		w.WriteHeader(c)
		return
	}

	if c != http.StatusOK {
		w.WriteHeader(c)
		return
	}

	name := strings.Split(r.URL.EscapedPath(), "/")[2] // todo 2

	e = es.AddVersion(name, hash, size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(c)
}

func get(w http.ResponseWriter, r *http.Request) {
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

	stream, err := getStream(object)
	if err != nil {
		utils.Log.Println(utils.Info, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	io.Copy(w, stream)
}

func getStream(object string) (io.Reader, error) {
	server := locate.Locate(object)
	if server == "" {
		return nil, fmt.Errorf("object %s locate fail", object)
	}
	return objectstream.NewGetStream(server, object)
}

func storeObject(r io.Reader, hash string, size int64) (int, error) {
	if locate.Exist(url.PathEscape(hash)) {
		utils.Log.Println(utils.Debug, hash, "file exist")
		// 存在就直接返回 ok
		return http.StatusOK, nil
	}

	stream, e := putStream(url.PathEscape(hash), size)
	if e != nil {
		return http.StatusInternalServerError, e
	}
	reader := io.TeeReader(r, stream)			// TeeReader 可以实现在读取的同事，将内容写入stream
	d := utils.CalculateHash(reader)
	if d != hash {
		// 验证hash值不通过
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch, calculated=%s, requestd=%s", d, hash)
	}

	stream.Commit(true)
	return http.StatusOK, nil
}

func putStream(hash string, size int64) (*objectstream.TempPutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}
	utils.Log.Println(utils.Debug, "select data server ", server)
	return objectstream.NewTempPutStream(server, hash, size)
}
