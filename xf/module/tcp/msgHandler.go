package tcp

import (
	"lib/Message"
	"lib/glog"
	"lib/link"
	"lib/proto"
	"xf/module"
)

func (self *TcpModule) msgHandler(rawMsg []byte, ch chan *link.Message, sid uint64) {
	pbmsg := new(Message.Message)

	err := proto.Unmarshal(rawMsg, pbmsg)
	if err != nil {
		glog.Errorf("TCP: msgHandler err [%s]", err)
		return
	}
	chRec := make(chan *map[uint64]*Message.Message)

	go module.Msg.Handle(pbmsg, sid, chRec)

	for pmsg := range chRec {
		if pmsg != nil { // msg nil means no rec msgHandler
			for tosid, msg := range *pmsg {
				bymsg, err := proto.Marshal(msg)
				if err != nil {
					glog.Errorf("msgHandler: proto.Marshal err [%s]", err.Error())
					continue
				}
				linkmsg := link.Message(link.Binary(bymsg))
				if tosid == sid {
					ch <- &linkmsg
				} else {
					//err := self.server.TrySend(tosid, linkmsg, time.Second)
					go self.server.Send(tosid, linkmsg)
				}
			}
		}
	}

	close(ch)
}
