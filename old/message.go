package cchelper

import (
	"fmt"
	"lib/Message"
	"lib/proto"
)

// region errhandler
var (
// 定义错误
)

var mapReq2Res map[Message.MsgType]Message.MsgType = map[Message.MsgType]Message.MsgType{
	Message.MsgType_Login_Request:      Message.MsgType_Login_Response,
	Message.MsgType_Help_Request:       Message.MsgType_Help_Response,
	Message.MsgType_Chat_Request:       Message.MsgType_Chat_Response,
	Message.MsgType_Broadcast_Request:  Message.MsgType_Broadcast_Response,
	Message.MsgType_UserManage_Request: Message.MsgType_UserManage_Response,
}

var errResponse = func(msg *Message.Message, errMsg string) *Message.Message {
	return &Message.Message{
		Type:     mapReq2Res[msg.GetType()].Enum(),
		Sequence: proto.Uint32(msg.GetSequence()),
		Response: &Message.Response{
			Result:           proto.Bool(false),
			ErrorDescription: []byte(errMsg),
		},
	}
}

// endregion errorHandler

// region messageLoop
type messageHandler func(*Message.Message, int64) *Message.Message

// todo 此处map可以优化成为slice Message.MsgType 为slice 索引
var mapMsgHandler = map[Message.MsgType][]messageHandler{
	Message.MsgType_Keepalive: []messageHandler{
		func(msg *Message.Message, uuid int64) *Message.Message {
			//fmt.Printf("MSG_Keepalive (%s)\n", msg.String())
			return &Message.Message{
				Type:     Message.MsgType_Keepalive.Enum(),
				Sequence: proto.Uint32(msg.GetSequence()),
			}
		},
	},
	Message.MsgType_Login_Request: []messageHandler{
		func(msg *Message.Message, uuid int64) *Message.Message {
			//fmt.Printf("MSG_Login_Request (%s)\n", msg.String())
			request := msg.GetRequest()
			if request == nil {
				return errResponse(msg, "")
			}
			login := request.GetLogin()
			if login == nil {
				return errResponse(msg, "")
			}

			onl := login.GetOnline()
			if onl != false {
				name := string(login.GetUsername())
				//pwd := login.GetPassword()
				//visi := login.GetVisibility()

				//Clients.rwm.Lock()
				Clients.Uuid2c[uuid].username = name
				Clients.Name2c[name] = Clients.Uuid2c[uuid]
				//Clients.rwm.Unlock()

			} else {
				// 登出处理... 貌似不用处理...因为登出必然会断开连接..断开连接会自动处理

			}
			return &Message.Message{
				Type:     Message.MsgType_Login_Response.Enum(),
				Sequence: proto.Uint32(msg.GetSequence()),
				Response: &Message.Response{
					Result: proto.Bool(true),
					//ErrorDescription: []byte("Hello Error"),
				},
			}
		},
	},
	Message.MsgType_Help_Request: []messageHandler{
		func(msg *Message.Message, uuid int64) *Message.Message {
			//fmt.Printf("MsgType_Help_Request (%s)\n", msg.String())
			return &Message.Message{
				Type:     Message.MsgType_Help_Response.Enum(),
				Sequence: proto.Uint32(msg.GetSequence()),
				Response: &Message.Response{
					Result: proto.Bool(true),
					Help: &Message.HelpResponse{
						Helper: &Message.UserInfo{
							From:     []byte("helper"),
							Realname: proto.String("Realname"),
							Location: proto.String("Location"),
						},
					},
				},
			}
		},
		// help 统计数据
		func(msg *Message.Message, uuid int64) *Message.Message {
			// todo 通过chanel 改写tcp_statis中help统计数据
			chHelp <- uuid
			chPublic <- uuid
			return nil
		},
	},
	Message.MsgType_Chat_Request: []messageHandler{
		func(msg *Message.Message, uuid int64) *Message.Message {
			fmt.Printf("MsgType_Chat_Request (%s)\n", msg.String())
			return &Message.Message{
				Type:     Message.MsgType_Chat_Response.Enum(),
				Sequence: proto.Uint32(msg.GetSequence()),
				Response: &Message.Response{
				//Result:           proto.Bool(false),
				//ErrorDescription: []byte("Hello Error"),
				},
			}
		},
	},
	Message.MsgType_Broadcast_Request: []messageHandler{
		func(msg *Message.Message, uuid int64) *Message.Message {
			fmt.Printf("MsgType_Broadcast_Request (%s)\n", msg.String())
			return &Message.Message{
				Type:     Message.MsgType_Broadcast_Response.Enum(),
				Sequence: proto.Uint32(msg.GetSequence()),
				Response: &Message.Response{
				//Result:           proto.Bool(false),
				//ErrorDescription: []byte("Hello Error"),
				},
			}
		},
	},
	Message.MsgType_UserManage_Request: []messageHandler{
		func(msg *Message.Message, uuid int64) *Message.Message {
			fmt.Printf("MsgType_UserManage_Request (%s)\n", msg.String())
			return &Message.Message{
				Type:     Message.MsgType_UserManage_Response.Enum(),
				Sequence: proto.Uint32(msg.GetSequence()),
				Response: &Message.Response{
				//Result:           proto.Bool(false),
				//ErrorDescription: []byte("Hello Error"),
				},
			}
		},
	},
}

type messageServer struct {
	ctx *CCHelper
}

func (ms *messageServer) Handler(msgReq *chanMsgReq) {
	//glog.V(3).Infof("Message Handler processes message (%s)", msgReq.msg.String())

	if handlers, ok := mapMsgHandler[msgReq.msg.GetType()]; ok {
		// 收到一条可回复n条
		for _, handler := range handlers {
			// 如果此处改写成goroutine 匿名函数,必须建立一个chan ,通知关闭res
			res := handler(msgReq.msg /*, msgReq.clientConn*/, msgReq.uuid)
			// 返回nil 表示handler 为注册的某种类型消息处理插件
			if res != nil {
				msgReq.res <- res
			}
		}
		close(msgReq.res)
	}
}
