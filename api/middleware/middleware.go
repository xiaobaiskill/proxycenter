package middleWare

import (
	"net/http"
)

type middleWare struct {
	r *http.ServeMux
}

func (m middleWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.r.ServeHTTP(w, r)
}

func NewMiddleWare(mux *http.ServeMux) middleWare {
	m := middleWare{}
	m.r = mux
	return m
}
