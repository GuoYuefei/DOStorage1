package temp

import (
	"fmt"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/rs"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"net/http"
	"strings"
)

func head(w http.ResponseWriter, r *http.Request) {
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
	utils.Log.Printf(utils.Debug, "get current len is %v\n", current)
	w.Header().Set("content-length", fmt.Sprintf("%d", current))
}
