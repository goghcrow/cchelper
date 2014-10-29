package tcp

import (
	"lib/glog"
	"lib/link"
	"time"
)

func (self *TcpModule) sessionHanlder(session *link.Session) {
	glog.V(3).Infof("TCP: session id [%d] start ", session.Id())

	for {
		rawMsg, err := session.Read()
		if err != nil {
			glog.Errorf("TCP: seesionHandler read err [%s]", err.Error())
			session.Close(err)
			break
		}

		// 心跳超时
		err = session.Conn().SetReadDeadline(time.Now().Add(self.readOverTime))
		if err != nil {
			glog.Errorf("TCP: seesionHandler [id:%d] SetReadDeadline Fail", session.Id())
		}

		chRec := make(chan *link.Message)

		go self.msgHandler(rawMsg, chRec, session.Id())

		for msg := range chRec {
			//err = session.TrySend(msg, time.Second*2) // async
			//if err != nil {
			//	fmt.Println(err)
			//}
			session.Send(*msg) // sync
		}
	}

	glog.V(3).Infof("TCP: session id [%d] close  ", session.Id())
}
