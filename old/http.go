package cchelper

import (
	"net/http"
	"net/http/pprof"
)

type httpServer struct {
	ctx *CCHelper
}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	pprof.Index(w, req)
}
