package versions

import (
	"encoding/json"
	"github.com/GuoYuefei/DOStorage1/distributed/es"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	from := 0
	size := 1000
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	for {
		metas, e := es.SearchAllVersions(name, from, size)
		if e != nil {
			utils.Log.Println(utils.Err, e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		utils.Log.Println(utils.Debug, metas)
		for i := range metas {
			b, _ := json.Marshal(metas[i])
			w.Write(b)
			w.Write([]byte{'\n'})
		}
		if len(metas) != size {
			return
		}
		from += size
	}
}




