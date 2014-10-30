package tcp

import (
	"encoding/binary"
	"lib/Message"
	"lib/glog"
	"lib/link"
	"lib/proto"
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
	channels     map[uint64]*link.Channel // 群组,0为全局群组
}

func init() {
	module.Tcp = &TcpModule{
		channels: make(map[uint64]*link.Channel),
	}
}

func (self *TcpModule) InitChannel() {

}

// 开启tcp服务,实现module接口
func (self *TcpModule) Start(listener net.Listener) {

	// 服务器初始化
	self.readOverTime = time.Duration(int64(time.Second) * int64(module.Opt.ReadTOver())) // 初始化心跳超时时间
	self.tcpMaxConn = module.Opt.MaxConn()                                                // 初始化服务器接受最大维持长连接数
	self.readBuffer = module.Opt.ReadBuffer()
	self.protocol = link.PacketN(4, binary.BigEndian)     // 初始化包头4byte包体长度二进制协议
	self.server = link.NewServer(listener, self.protocol) // 初始化server实例
	self.server.SetReadBufferSize(self.readBuffer)
	sem := limit.NewMaxBlockCtrl(self.tcpMaxConn)              // 初始化长连接控制器
	self.channels[0] = link.NewChannel(self.server.Protocol()) // 0 global channel

	// 注册异常处理
	defer func() {
		glog.Infof("TCP: closing %s", listener.Addr())
		if err := recover(); err != nil {
			glog.Fatalf("TCP: server panic(%s)", err)
			if !self.server.IsStopped() {
				self.server.Stop(err)
			}
		}
	}()

	// 测试全局广播
	/*
		go func() {
			for {
				time.Sleep(time.Second)
				msg := module.Msg.GenKeepaliveResp(0)
				bymsg, err := proto.Marshal(msg)
				if err != nil {
					glog.Errorf("msgHandler: proto.Marshal err [%s]", err.Error())
					continue
				}
				self.channels[0].Broadcast(link.Binary(bymsg))
			}
		}()
	*/

	glog.Infof("TCP: listening on %s", listener.Addr())

	// tcp server 大循环
	for {
		sem.Pop() // 最大长连接数控制,弹出

		session, err := self.server.Accept()
		// 防止tcp退出时,session关闭之后等于nil
		if session != nil {
			glog.V(3).Infof("TCP: client %s in", session.Conn().RemoteAddr().String())
		}
		if err != nil {
			glog.Errorf("TCP: server accept err [%s]", err.Error())
			self.server.Stop(err) // todo: 接收异常关闭服务器 ?
			break
		}

		module.Sta.Online() // 上线:人数统计

		// 非阻塞处理每一个session
		go func(session *link.Session) { // 传参复制session指针,保持引用
			defer func() {
				sem.Push()           // 最大长连接数控制,压入
				module.Sta.Offline() // 离线:人数统计
				glog.V(3).Infof("TCP: client %s close", session.Conn().RemoteAddr().String())
			}()

			self.channels[0].Join(session, nil) // 加入全局channel,无踢出回调
			//glog.V(3).Infof("TCP: session id [%d] start ", session.Id())

			// 注册异步发送(trySend)失败处理函数
			// todo 发送失败处理...
			session.OnSendFailed = func(session *link.Session, err error) {
				if session != nil {
					glog.Errorf("TCP: sid [%d] TrySend failed [%s]", session.Id(), err.Error())
				}
			}

			// session处理循环
			for {
				rawMsg, err := session.Read()
				if err != nil {
					glog.Errorf("TCP: seesionHandler read err [%s]", err.Error())
					session.Close(err) // todo: 接收失败关闭session ?
					break
				}
				err = session.Conn().SetReadDeadline(time.Now().Add(self.readOverTime)) // 心跳超时
				if err != nil {
					glog.Errorf("TCP: seesionHandler [id:%d] SetReadDeadline Fail", session.Id())
					session.Close(err) // todo: 心跳超时失败关闭session ?
					break
				}

				// message处理
				// todo: 暂时放弃goroutine处理,因为客户端发送连续消息会导致粘包,解析错误
				//go func(rawMsg []byte) {
				pbmsg := new(Message.Message)
				err = proto.Unmarshal(rawMsg, pbmsg)
				if err != nil {
					glog.Errorf("TCP: msgHandler Unmarshal err [%s]", err)
					return
				}

				// todo: 将map置换为struct
				//chRec := make(chan *map[uint64]*Message.Message)
				chRec := make(chan *[]*module.MsgPack)
				go module.Msg.Handle(pbmsg, session, chRec)

				for msgPackes := range chRec {
					if msgPackes != nil { // 消息为nil,表示为非回复消息处理返回值
						//for tosid, msg := range *msgPack {
						for _, msgPack := range *msgPackes {
							//bymsg, err := proto.Marshal(msg)
							bymsg, err := proto.Marshal(msgPack.Msg)
							if err != nil {
								glog.Errorf("TCP: msgHandler Marshal err [%s]", err.Error())
								continue
							}
							// todo: 异步发送消息, 超时0,默认缓存1024条消息,超过触发异常 ? 适当可以调节增大channel缓存大小
							re := 1 // 重发计数
						trysend:
							//err = self.server.TrySend(tosid, link.Binary(bymsg), 0)
							err = self.server.TrySend(msgPack.Sid, link.Binary(bymsg), 0)
							switch err {
							case link.BlockingError:
								// todo: 是否存在会发生死循环呢?? 这里需要重点关注下
								// 发生阻塞,自动设置二倍sendchan大小,重发
								newSendChanSize := self.server.GetSendChanSize() << 1
								if newSendChanSize <= 4096 { // sendChanSize 过大会产生发送延迟
									self.server.SetSendChanSize(newSendChanSize)
									glog.Infof("TCP: SendChanSize expand to %d", newSendChanSize)
								}
								re++
								if re < 4 {
									goto trysend
								} else {
									// 出现此错误说明服务器负载太大了.....
									glog.Errorf("TCP: try send failed 3 times")
								}
							case nil:
							case link.SendToClosedError:
								glog.Error("TCP: try send to closed session")
							default: // 其他错误
								glog.Errorf("TCP: try send error [%s]", err.Error())
							}
							// 同步发送
							// 如果多条消息发送给一个session,没有必要使用goroutine,send会对session上锁
							// 如果多条消息发送个不同session,可以使用goroutine
							// 返回error未处理
							//go self.server.Send(tosid, link.Binary(bymsg))
						}
					}
				}
				//}(rawMsg)
			}
			//glog.V(3).Infof("TCP: session id [%d] close  ", session.Id())
		}(session)
	}
}

// 停止tcp服务,实现module接口
func (self *TcpModule) Stop(reasion interface{}) {
	self.server.Stop(reasion)
}
