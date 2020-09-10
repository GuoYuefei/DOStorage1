package objects

import (
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	utils.Log.Printf(utils.Debug, "will call function objects.%s", m)
	if m == http.MethodGet {
		get(w, r)
		return
	}

	if m == http.MethodPost {
		post(w, r)
		return
	}

	if m == http.MethodPut {
		put(w, r)
		return
	}

	if m == http.MethodDelete {
		del(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
