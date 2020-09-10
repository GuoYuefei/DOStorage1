package temp

import (
	"github.com/GuoYuefei/DOStorage1/distributed/config"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

func patch(w http.ResponseWriter, r *http.Request) {
	uu := strings.Split(r.URL.EscapedPath(), "/")[2]
	tempinfo, e := readFromFile(uu)
	if e != nil {
		utils.Log.Println(utils.Err, e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFile := path.Join(config.TempRoot, uu)
	datFile := infoFile+".dat"
	f, e := os.OpenFile(datFile, os.O_WRONLY|os.O_APPEND, 0)
	//f, e := os.Create(datFile)
	if e != nil {
		utils.Log.Println(utils.Err, e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, e = io.Copy(f, r.Body)
	if e != nil {
		utils.Log.Println(utils.Err, e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	info, e := f.Stat()
	if e != nil {
		utils.Log.Println(utils.Err, e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	actual := info.Size()
	if actual > tempinfo.Size {
		os.Remove(datFile)
		os.Remove(infoFile)
		utils.Log.Printf(utils.Debug, "actual size %d, exceeds %d", actual, tempinfo.Size)
		w.WriteHeader(http.StatusInternalServerError)
	}
}