package temp

import (
	"github.com/GuoYuefei/DOStorage1/distributed/es"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/locate"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/rs"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/util"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, e := rs.NewRSResumablePutStreamFromToken(token)
	if e != nil {
		utils.Log.Println(utils.Err, e)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	current := stream.CurrentSize()
	if current == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	offset := util.GetOffsetFromHeader(r.Header)
	if current != offset {
		// 现在的大小和offset不匹配
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}
	bytes := make([]byte, rs.BLOCK_SIZE)
	for {
		n, e := io.ReadFull(r.Body, bytes)
		if e != nil && e != io.EOF && e != io.ErrUnexpectedEOF {
			utils.Log.Println(utils.Err, e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		current += int64(n)
		if current > stream.Size {
			stream.Commit(false)
			utils.Log.Println(utils.Err, "resumable put exceed size")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		utils.Log.Printf(utils.Debug, "temp.put size: %v \t current %v \t stream.size %v\n", n, current, stream.Size)
		// 上传出现问题， 最后一次没达到block_size 大小，也不是最后一块， 直接丢弃
		if n != rs.BLOCK_SIZE && current != stream.Size {
			return
		}
		stream.Write(bytes[:n])
		if current == stream.Size {
			// 全部写入后的操作
			stream.Flush()
			getStream, e := rs.NewRSResumableGetStream(stream.Servers, stream.Uuids, stream.Size)
			if e != nil {
				stream.Commit(false)
				utils.Log.Println(utils.Err, e)
				w.WriteHeader(http.StatusForbidden)
				return
			}
			hash := utils.CalculateHash(getStream)
			if hash != stream.Hash {
				stream.Commit(false)
				utils.Log.Println(utils.Err, "resumable put done but hash mismatch")
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if locate.Exist(url.PathEscape(hash)) {
				utils.Log.Println(utils.Debug, "object is exist before")
				stream.Commit(false)
			} else {
				stream.Commit(true)
			}
			utils.Log.Println(utils.Debug, "suss upload!!!")
			e = es.AddVersion(stream.Name, stream.Hash, stream.Size)
			if e != nil {
				utils.Log.Println(utils.Err, e)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	}
}
