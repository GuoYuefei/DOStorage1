package objects

import (
	"github.com/GuoYuefei/DOStorage1/distributed/data/locate"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	utils.Log.Println(utils.Debug, "objects Information received is ", m)
	if m == http.MethodGet {
		get(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func get(w http.ResponseWriter, r *http.Request) {
	file := getFile(strings.Split(r.URL.EscapedPath(),"/")[2])
	if file == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err := sendFile(w, file)
	if err != nil {
		utils.Log.Println(utils.Err, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getFile(hash string) string {
	file := path.Join(locate.ObjectRoot, hash)
	f, e := os.Open(file)
	if e != nil {
		return ""
	}
	d := url.PathEscape(utils.CalculateHash(f))
	f.Close()
	if d != hash {
		// 验证文件不成功
		utils.Log.Printf(utils.Err, "object hash mismatch, remove", file)
		locate.Del(hash)
		os.Remove(file)
		return ""
	}
	return file
}

func sendFile(writer io.Writer, file string) error {
	f, e := os.Open(file)
	if e != nil {
		return e
	}
	defer f.Close()
	io.Copy(writer, f)
	return nil
}
