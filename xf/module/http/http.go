package http

import (
	"lib/glog"
	"net"
	"net/http"
	"strings"
	"xf/module"
)

type HttpModule struct {
	server *http.Server
}

func init() {
	module.Http = &HttpModule{}
}

func (self *HttpModule) Start(listener net.Listener) {
	glog.Infof("HTTP: listening on %s", listener.Addr())

	defer func() {
		if err := recover(); err != nil {
			glog.Fatalf("HTTP: server panic(%s)", err)

			//glog.Info("HTTP: server stop")
			//self.Stop(err)
		}
	}()

	self.server = &http.Server{Handler: httpServer{}}
	err := self.server.Serve(listener)
	// theres no direct way to detect this error because it is not exposed
	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		glog.Errorf("ERROR: http.Serve() - %s", err)
	}

	glog.Infof("HTTP: closing %s", listener.Addr())
}

// todo 实现stop方法清理http资源
func (self *HttpModule) Stop(reasion interface{}) {

}
