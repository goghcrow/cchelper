package cchelper

import (
	"lib/glog"
	"net"
	"runtime"
	"strings"
)

type TCPHandler interface {
	Handler(net.Conn)
}

func TCPServer(listener net.Listener, handler TCPHandler) {
	glog.Infof("TCP: listening on %s", listener.Addr())

	defer func() {
		if e := recover(); e != nil {
			glog.Fatalf("TCP: server panic(%s)", e)
		}
	}()

	for {
		clientConn, err := listener.Accept()

		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				glog.Infof("temporary Accept() failure - %s", err)
				runtime.Gosched() // 让出cpu片段
				continue
			}

			// theres no direct way to detect this error because it is not exposed
			if !strings.Contains(err.Error(), "use of closed network connection") {
				glog.Errorf("listener.Accept() - %s", err)
			}
			break
		}

		go handler.Handler(clientConn)

		// clientConn.SetReadDeadline()
		// 加入一个连接检测机制
		//go func(conn net.Conn) {
		//w := timingwheel.NewTimingWheel(time.Second, 60) // 1分钟检查一次连接

		//for {
		//	select {
		//	case <-w.After(time.Second):
		//		if err := ping(); err != nil {
		//			conn.Close()
		//		}
		//		return
		//	}
		//}
		//}(clientConn)

	}

	glog.Infof("TCP: closing %s", listener.Addr())
}
