package objects

import (
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodGet {
		get(w, r)
		return
	}

	if m == http.MethodPut {
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
