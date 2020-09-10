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

func get(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	f, e := os.Open(path.Join(config.ServerData.STORAGE_ROOT, "temp", uuid+".dat"))
	if e != nil {
		utils.Log.Println(utils.Err, e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	io.Copy(w, f)
}

