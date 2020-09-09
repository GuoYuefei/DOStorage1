package temp

import (
	"encoding/json"
	"github.com/GuoYuefei/DOStorage1/distributed/config"
	"github.com/GuoYuefei/DOStorage1/distributed/data/locate"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	utils.Log.Printf(utils.Debug, "will call function temp.%s", m)
	switch m {
	case http.MethodHead: head(w, r)
	case http.MethodGet: get(w, r)
	case http.MethodPut: put(w, r)
	case http.MethodPatch: patch(w, r)
	case http.MethodPost: post(w, r)
	case http.MethodDelete: delete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func readFromFile(uu string) (*tempInfo, error) {
	f, e := os.Open(path.Join(config.TempRoot, uu))
	if e != nil {
		return nil, e
	}
	defer f.Close()
	b, _ := ioutil.ReadAll(f)
	var info tempInfo
	json.Unmarshal(b, &info)
	return &info, nil
}

// 最后的名字定位成 <hash>.X.<hash of shard X>
// todo 错误处理
func commitTempObject(datFile string, info *tempInfo) {
	f, e := os.Open(datFile)
	if e != nil {
		utils.Log.Println(utils.Exception, e)
	}
	d := url.PathEscape(utils.CalculateHash(f))
	e = f.Close()
	if e != nil {
		time.Sleep(50*time.Millisecond)
		f.Close()
	}
	// 最多花费1秒
	for i := 10; i > 0; i-- {
		e = os.Rename(datFile, path.Join(config.ObjectRoot, info.Name+"."+d))
		if e != nil {
			time.Sleep(100*time.Millisecond)
		} else {
			break
		}
	}

	locate.Add(info.hash(), info.id())
}
