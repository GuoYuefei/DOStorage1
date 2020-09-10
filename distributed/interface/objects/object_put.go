package objects

import (
	"fmt"
	"github.com/GuoYuefei/DOStorage1/distributed/es"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/heartbeat"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/locate"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/rs"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/util"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

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

// 接收的 hash 是已经 url.pathescape 过得
func storeObject(r io.Reader, hash string, size int64) (int, error) {
	if locate.Exist(hash) {
		utils.Log.Println(utils.Debug, hash, "file exist")
		// 存在就直接返回 ok
		return http.StatusOK, nil
	}

	stream, e := putStream(hash, size)
	if e != nil {
		return http.StatusInternalServerError, e
	}
	reader := io.TeeReader(r, stream)			// TeeReader 可以实现在读取的同事，将内容写入stream
	d := utils.CalculateHash(reader)
	d = url.PathEscape(d)
	if d != hash {
		// 验证hash值不通过
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch, calculated=%s, requestd=%s", d, hash)
	}

	stream.Commit(true)
	return http.StatusOK, nil
}

func putStream(hash string, size int64) (*rs.RSPutStream, error) {
	servers := heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS, nil)
	if len(servers) != rs.ALL_SHARDS {
		return nil, fmt.Errorf("cannot find enough dataServer")
	}

	return rs.NewRSPutStream(servers, hash, size)
}
