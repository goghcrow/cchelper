// 84 105 125 有问题
// 大并发异常退出 这些chan有问题
package cchelper

import (
	"io"
	"lib/Message"
	"lib/glog"
	"lib/proto"
	"net"
	"time"
)

// region tcpServer
type tcpServer struct {
	ctx *CCHelper
}

type chanReq struct {
	req net.Conn
	res chan *[]byte
}

// 2014-10-16 add outstanding contronl 10000
const MaxOutstanding = 10000

type null struct{}

var sem = make(chan null, MaxOutstanding)

func init() {
	// outstanding init
	for i := 0; i < MaxOutstanding; i++ {
		sem <- null{}
	}
}

// todo 处理各种tcp异常关闭情况 to see @from http://blog.csdn.net/icyday/article/details/20638321
func (s *tcpServer) Handler(clientConn net.Conn) {
	uuid := Uuid()

	defer func() {
		if err := recover(); err != nil {
			glog.Fatalf("TCP Handler panic(%s)", err)
		}
	}()

	go func() {
		chOnline <- 1 // online + 1
	}()

	// 假设接到的客户端第一条消息,消息头正确

	// packet v2
	//chFromUnpack := make(chan *Packet) // 读取解包数据
	//chFromMarMsg := make(chan *Packet) // 读取序列化之后消息

	// packet v1
	chFromUnpack := make(chan []byte) // 读取解包数据
	chFromMarMsg := make(chan []byte) // 读取序列化之后消息 // 通过此chan发送byte[]给客户端

	// 是否应该加锁呢???
	//Clients.rwm.RLock()
	Clients.Uuid2c[uuid] = &client{chFromMarMsg: chFromMarMsg}
	//Clients.rwm.RUnlock()

	// 2014-10-16 最大长连接数控制
	<-sem
	defer func() {
		sem <- null{}

		//_, ok := <-chFromMarMsg // 这里会引起阻塞,不能判断
		//if ok {
		close(chFromMarMsg)
		//}
		//_, ok = <-chFromUnpack // 这里会引起阻塞,不能判断
		//if ok {
		close(chFromUnpack)
		//}

		clientConn.Close()
		//clientConn = nil

		//Clients.rwm.RLock()
		if client, ok := Clients.Uuid2c[uuid]; ok {
			delete(Clients.Name2c, client.username)
		}
		delete(Clients.Uuid2c, uuid)
		//Clients.rwm.RUnlock()

		chOnline <- -1 // online - 1
	}()

	// protobuf 回复客户端
	go func(uuid int64) {
		for {
			//debugRes++
			//fmt.Println("debugRes ", debugRes)
			packet, ok := <-chFromMarMsg // 同步阻塞
			if !ok {
				return
			}
			//fmt.Println("<-chFromMarMsg : ", packet)
			glog.V(3).Infof("TCP Handler pack&send msg to client[%s]", clientConn.RemoteAddr())

			// !!!!!!!!!!!!!!!
			// 此处需要在goroutine中写入...要不客户端发送过来的连续信息,只能接收到763条
			// 到底是为什么
			go clientConn.Write(pack(packet))
		}
	}(uuid)

	// 处理解包数据
	go func(uuid int64) {
		for {
			//debugPac++
			//fmt.Println("debugPac ", debugPac)

			// packet v1
			packet, ok := <-chFromUnpack // 同步阻塞,等待解包消息体
			if !ok {
				return
			}
			//fmt.Println("<-chFromUnpack : ", packet)

			msg := new(Message.Message)
			err := proto.Unmarshal(packet, msg)
			if err != nil {
				glog.Warningf("TCP: Handler Unmarshal Message err : %s", err)
				// !!!!!!!!通知客户端 消息包解序列化错误...
				continue
			}

			//glog.V(3).Infof("TCP Handler send msg (%s) to MessageLoop", msg.String())

			res := make(chan *Message.Message)
			s.ctx.tcpMsgChan <- &chanMsgReq{msg: msg /*, clientConn: &clientConn*/, uuid: uuid, res: res}

			// 解决一次性回复client n 条消息
			for recMsg := range res {
				//recMsg := <-res // 同步阻塞,等待回复消息,Message处理中可以提前返回回复消息,耗时操作go func()

				//glog.V(3).Infof("TCP Handler get msg (%s) from MessageLoop", recMsg.String())

				resbuf, err := proto.Marshal(recMsg)

				if err != nil {
					glog.Warningln("TCP: Handler Marshal Message err : ", err)
					// !!!!!!!!通知客户端 消息包序列化错误...
					continue
				}
				chFromMarMsg <- resbuf
			}

			// packet v2
			/*
				resPac := &Packet{}

				packet, ok := <-chFromUnpack // 同步阻塞,等待解包消息体
				if !ok {
					return
				}
						// 添加新协议时候在此添加...
						switch packet.ptype {
						case Pac_Type_HeartBeat:
							resPac.ptype = Pac_Type_HeartBeat
							//resPac.pdata = tmp
						case Pac_Type_Protobuf:
							msg := new(Message.Message)
							err := proto.Unmarshal(packet.pdata, msg)
							if err != nil {
								glog.Warningf("TCP: Handler Unmarshal Message err : %s", err)
								// !!!!!!!!通知客户端 消息包解序列化错误...
								continue
							}

							//glog.V(3).Infof("TCP Handler send msg (%s) to MessageLoop", msg.String())

							res := make(chan *Message.Message)
							s.ctx.tcpMsgChan <- &chanMsgReq{msg: msg, res: res}
							recMsg := <-res // 同步阻塞,等待回复消息,Message处理中可以提前返回回复消息,耗时操作go func()

							//glog.V(3).Infof("TCP Handler get msg (%s) from MessageLoop", recMsg.String())

							resbuf, err := proto.Marshal(recMsg)

							if err != nil {
								glog.Warningln("TCP: Handler Marshal Message err : ")
								// !!!!!!!!通知客户端 消息包序列化错误...
								continue
							}

							resPac.ptype = Pac_Type_Protobuf
							resPac.pdata = resbuf
						default:
							// 未识别消息类型返回给客户端
							resPac.ptype = Pac_Type_Unknown
							//resPac.pdata =
						}

					chFromMarMsg <- resPac
			*/

		}
	}(uuid)

	tmpBuffer := make([]byte, 0) // 被截断数据
	buffer := make([]byte, 1024) // 接受客户端数据,分包

	// 从无限等待变为按心跳包时间超时
	// 读取无超时...保持连接...需要在TCPServer中添加一个连接是否正常的保持机制
	//	clientConn.SetReadDeadline(time.Time{})

	for {
		// 2014-10-19 添加心跳包超时
		err := clientConn.SetReadDeadline(time.Now().Add(s.ctx.opts.TCPReadTimeout))
		if err != nil {
			glog.Fatalln("TCP: Handler [", clientConn.RemoteAddr(), "] SetReadDeadline fail :", err)
			// 这里可能导致客户端永远在线...
		}

		n, err := clientConn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				glog.Infof("TCP: Handler detected client (%s) closed connection", clientConn.RemoteAddr())
			} else {
				glog.Infof("TCP: Handler clientConn.Read() err: %s", err.Error())
			}

			// todo 断线的清理工作....

			// 转移到defer
			////Clients.rwm.RLock()
			//if client, ok := Clients.Uuid2c[uuid]; ok {
			//	delete(Clients.Name2c, client.username)
			//}
			//delete(Clients.Uuid2c, uuid)
			////Clients.rwm.RUnlock()

			// 转移到defer
			//close(chFromUnpack) // => panic: runtime error: send on closed channel
			//close(chFromMarMsg) // => panic: runtime error: send on closed channel

			// clientConn.Close() // => panic: runtime error: invalid memory address or nil pointer dereference
			// clientConn = nil
			break
		}

		//fmt.Println("clientConn.Read(buffer) : ", buffer[:n])
		tmpBuffer = unpack(append(tmpBuffer, buffer[:n]...), chFromUnpack)
	}
}
