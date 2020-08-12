package util

import (
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"net/http"
	"strconv"
)

func GetHashFromHeader(h http.Header) string {
	digest := h.Get("digest")
	if len(digest) < 9 {
		return ""
	}
	if digest[:8] != "SHA-256=" {
		return ""
	}
	return digest[8:]
}

func GetSizeFromHeader(h http.Header) int64 {
	size, e := strconv.ParseInt(h.Get("Content-Length"), 0, 64)
	if e != nil {
		utils.Log.Println(utils.Err, e)
	}
	return size
}
