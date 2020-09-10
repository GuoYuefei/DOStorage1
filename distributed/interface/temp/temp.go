package temp

import (
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	utils.Log.Printf(utils.Debug, "will call function temp.%s", m)
	if m == http.MethodPut {
		put(w, r)
		return
	}

	if m == http.MethodHead {
		head(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}