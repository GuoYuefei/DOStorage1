package rs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/objectstream"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/util"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"io"
	"net/http"
)

type resumableToken struct {
	Name string
	Size int64
	Hash string
	Servers []string
	Uuids []string
}

type RSResumablePutStream struct {
	*RSPutStream
	*resumableToken
}

func NewRSResumablePutStreamFromToken(token string) (*RSResumablePutStream, error) {
	// TODO 对应 ToToken 这个方法
	b, e := base64.StdEncoding.DecodeString(token)
	if e != nil {
		return nil, e
	}

	var t resumableToken
	e = json.Unmarshal(b, &t)
	if e != nil {
		return nil, e
	}

	writers := make([]io.Writer, ALL_SHARDS)
	for i := range writers {
		writers[i] = &objectstream.TempPutStream{t.Servers[i], t.Uuids[i]}
	}
	enc := NewEncoder(writers)
	return &RSResumablePutStream{&RSPutStream{enc}, &t}, nil
}

// 返回第一个临时分片的大小*4， 若超出文件大小，则返回文件大小
func (s *RSResumablePutStream) CurrentSize() int64 {
	r, e := http.Head(fmt.Sprintf("http://%s/temp/%s", s.Servers[0], s.Uuids[0]))
	if e != nil {
		utils.Log.Println(utils.Err, e)
		return -1
	}
	if r.StatusCode != http.StatusOK {
		utils.Log.Println(utils.Err, "RSResumablePutStream CurrentSize response is not OK, ")
		return -1
	}
	size := util.GetOffsetFromHeader(r.Header) * DATA_SHARDS
	if size > s.Size {
		size = s.Size
	}
	return size
}

func NewRSResumablePutStream(dataServers []string, name, hash string, size int64) (*RSResumablePutStream, error) {
	putStream, e := NewRSPutStream(dataServers, hash, size)
	if e != nil {
		return nil, e
	}
	uuids := make([]string, ALL_SHARDS)
	for i := range uuids {
		uuids[i] = putStream.writers[i].(*objectstream.TempPutStream).Uuid
	}
	token := &resumableToken{
		Name:    name,
		Size:    size,
		Hash:    hash,
		Servers: dataServers,
		Uuids:   uuids,
	}
	return &RSResumablePutStream{putStream, token}, nil
}

func (s *RSResumablePutStream) ToToken() string {
	b, _ := json.Marshal(s)
	// 可直接使用URLEncoding todo 还需加密
	return base64.StdEncoding.EncodeToString(b)
}

