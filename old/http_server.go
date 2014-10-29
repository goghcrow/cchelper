package cchelper

import (
	"lib/glog"
	"net"
	"net/http"
	"strings"
)

func HTTPServer(listener net.Listener, handler http.Handler) {
	glog.Infof("HTTP: listening on %s", listener.Addr())

	defer func() {
		if err := recover(); err != nil {
			glog.Fatalf("HTTP: server panic(%s)", err)
		}
	}()

	server := &http.Server{Handler: handler}
	err := server.Serve(listener)
	// theres no direct way to detect this error because it is not exposed
	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		glog.Errorf("ERROR: http.Serve() - %s", err)
	}

	glog.Infof("HTTP: closing %s", listener.Addr())
}
