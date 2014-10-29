package http

import "net/http"

type httpServer struct{} // add Ctx,options

func (httpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//pprof.Index(w, req)
}
