package tcp

import (
	"encoding/binary"
	"lib/glog"
	"lib/link"
	"lib/util/limit"
	"net"
	"time"
	"xf/module"
)

type TcpModule struct {
	server       *link.Server
	protocol     link.PacketProtocol
	readBuffer   int
	readOverTime time.Duration
	tcpMaxConn   uint64
}

func init() {
	module.Tcp = &TcpModule{}
}

func (self *TcpModule) Start(listener net.Listener) {
	glog.Infof("TCP: listening on %s", listener.Addr())

	self.readOverTime = time.Duration(int64(time.Second) * int64(module.Opt.ReadTOver()))
	self.tcpMaxConn = module.Opt.MaxConn()
	self.readBuffer = module.Opt.ReadBuffer()

	sem := limit.NewMaxBlockCtrl(self.tcpMaxConn)

	// todo packet head 写入配置文件
	self.protocol = link.PacketN(4, binary.BigEndian)
	self.server = link.NewServer(listener, self.protocol)
	self.server.SetReadBufferSize(self.readBuffer)

	defer func() {
		if err := recover(); err != nil {
			glog.Fatalf("TCP: server panic(%s)", err)
			if !self.server.IsStopped() {
				self.server.Stop(err)
			}
		}
	}()

	for {
		sem.Pop()

		session, err := self.server.Accept()
		// session is nil after close
		if session != nil {
			glog.V(3).Infof("TCP: client %s in", session.Conn().RemoteAddr().String())
		}
		if err != nil {
			glog.Errorf("TCP: server accept err [%s]", err.Error())
			self.server.Stop(err)
			break
		}

		module.Sta.Online()

		go func() {
			defer func() {
				sem.Push()
				module.Sta.Offline()
				glog.V(3).Infof("TCP: client %s close", session.Conn().RemoteAddr().String())
			}()

			self.sessionHanlder(session)
		}()
	}

	glog.Infof("TCP: closing %s", listener.Addr())
}

func (self *TcpModule) Stop(reasion interface{}) {
	self.server.Stop(reasion)
}

func (self *TcpModule) Protocol() link.PacketProtocol {

}
