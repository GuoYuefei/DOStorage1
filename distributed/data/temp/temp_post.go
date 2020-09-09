package temp

import (
	"github.com/GuoYuefei/DOStorage1/distributed/config"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"github.com/google/uuid"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

func post(w http.ResponseWriter, r *http.Request) {
	uu := uuid.New().String()
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size, e := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if e != nil {
		utils.Log.Println(utils.Err, e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t := tempInfo{
		uu,
		name,
		size,
	}
	e = t.writeToFile()
	if e != nil {
		utils.Log.Println(utils.Err, e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, e = os.Create(path.Join(config.TempRoot, t.Uuid+".dat"))
	if e != nil {
		utils.Log.Println(utils.Err, e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(uu))
}