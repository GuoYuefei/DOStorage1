package temp

import (
	"fmt"
	"github.com/GuoYuefei/DOStorage1/distributed/config"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"net/http"
	"os"
	"path"
	"strings"
)

func head(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	// todo 可以选择不打开文件， 而直接stat
	f, e := os.Open(path.Join(config.ServerData.STORAGE_ROOT, "temp", uuid+".dat"))
	if e != nil {
		utils.Log.Println(utils.Err, e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	info, e := f.Stat()
	if e != nil {
		utils.Log.Println(utils.Err, e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-length", fmt.Sprintf("%d", info.Size()))
}
