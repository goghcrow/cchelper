package cchelper

import (
	"lib/Message"
	"lib/glog"
)

type chanMsgReq struct {
	msg *Message.Message
	//clientConn *net.Conn
	uuid int64
	res  chan *Message.Message
}

type MessageHandler interface {
	Handler(*chanMsgReq)
}

// 消息处理循环
func TCPMessageLoop(msgReq chan *chanMsgReq, handler MessageHandler) {
	defer func() {
		if err := recover(); err != nil {
			glog.Fatalf("TCP MessageLoop panic(%s)", err)
		}
	}()

	//for {
	//	req := <-msgReq
	//	glog.V(2).Infof("TCP<essageLoop: receive one Message ()", req.msg.String())
	//	go handler.Handler(req)
	//}

	// for-range 在chan 关闭后自动退出
	for req := range msgReq {
		// 2014-10-16 修改为闭包,使req独立于每个处理函数
		// ???? 需要这样做么???
		go func(req *chanMsgReq) {
			glog.V(2).Infof("TCPMessageLoop: receive one Message ()", req.msg.String())
			go handler.Handler(req)
		}(req)
	}
}
