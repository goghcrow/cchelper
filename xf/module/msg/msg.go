package msg

import (
	"fmt"
	"lib/Message"
	"lib/proto"
	"xf/module"
)

type MsgModule struct{}

func init() {
	module.Msg = &MsgModule{}
}

func (*MsgModule) Handle(msg *Message.Message, sid uint64, ch chan *map[uint64]*Message.Message) {
	if handlers, ok := mapMsgHandler[msg.GetType()]; ok {
		// 收到一条可回复n条
		for _, handler := range handlers {
			//func(){}()
			// 如果此处改写成goroutine 匿名函数,必须建立一个chan ,通知关闭res
			res := handler(msg, sid)

			// 返回nil 表示handler 为注册的某种类型消息处理插件
			if res != nil {
				ch <- res
			}
		}
	} else {
		// todo 整理错误信息
		errResponse(msg, NoHandler, sid)
	}

	close(ch)
}

type messageHandler func(*Message.Message, uint64) *map[uint64]*Message.Message

// todo 此处map可以优化成为slice Message.MsgType 为slice 索引
var mapMsgHandler = map[Message.MsgType][]messageHandler{
	Message.MsgType_Keepalive: []messageHandler{
		func(msg *Message.Message, sid uint64) *map[uint64]*Message.Message {
			//fmt.Printf("MSG_Keepalive (%s)\n", msg.String())
			return &map[uint64]*Message.Message{
				sid: &Message.Message{
					Type:     Message.MsgType_Keepalive.Enum(),
					Sequence: proto.Uint32(msg.GetSequence()),
				},
			}
		},
	},
	Message.MsgType_Login_Request: []messageHandler{
		func(msg *Message.Message, sid uint64) *map[uint64]*Message.Message {
			fmt.Printf("MSG_Login_Request (%s)\n", msg.String())
			request := msg.GetRequest()
			if request == nil {
				return errResponse(msg, "", sid)
			}
			login := request.GetLogin()
			if login == nil {
				return errResponse(msg, "", sid)
			}

			onl := login.GetOnline()
			if onl != false {
				//name := string(login.GetUsername())
				////pwd := login.GetPassword()
				////visi := login.GetVisibility()

				////Clients.rwm.Lock()
				//Clients.Uuid2c[sid].username = name
				//Clients.Name2c[name] = Clients.Uuid2c[sid]
				////Clients.rwm.Unlock()

			} else {
				// 登出处理... 貌似不用处理...因为登出必然会断开连接..断开连接会自动处理

			}
			return &map[uint64]*Message.Message{
				sid: &Message.Message{
					Type:     Message.MsgType_Login_Response.Enum(),
					Sequence: proto.Uint32(msg.GetSequence()),
					Response: &Message.Response{
						Result: proto.Bool(true),
						//ErrorDescription: []byte("Hello Error"),
					},
				},
			}
		},
	},
	Message.MsgType_Help_Request: []messageHandler{
		func(msg *Message.Message, sid uint64) *map[uint64]*Message.Message {
			//fmt.Printf("MsgType_Help_Request (%s)\n", msg.String())
			return &map[uint64]*Message.Message{
				sid: &Message.Message{
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
				},
			}
		},
		// help 统计数据
		func(msg *Message.Message, sid uint64) *map[uint64]*Message.Message {
			go module.Sta.Help(sid)
			go module.Sta.Public(sid)
			return nil
		},
	},
	Message.MsgType_Chat_Request: []messageHandler{
		func(msg *Message.Message, sid uint64) *map[uint64]*Message.Message {
			fmt.Printf("MsgType_Chat_Request (%s)\n", msg.String())
			return &map[uint64]*Message.Message{
				sid: &Message.Message{
					Type:     Message.MsgType_Chat_Response.Enum(),
					Sequence: proto.Uint32(msg.GetSequence()),
					Response: &Message.Response{
					//Result:           proto.Bool(false),
					//ErrorDescription: []byte("Hello Error"),
					},
				},
			}
		},
	},
	Message.MsgType_Broadcast_Request: []messageHandler{
		func(msg *Message.Message, sid uint64) *map[uint64]*Message.Message {
			fmt.Printf("MsgType_Broadcast_Request (%s)\n", msg.String())
			return &map[uint64]*Message.Message{
				sid: &Message.Message{
					Type:     Message.MsgType_Broadcast_Response.Enum(),
					Sequence: proto.Uint32(msg.GetSequence()),
					Response: &Message.Response{
					//Result:           proto.Bool(false),
					//ErrorDescription: []byte("Hello Error"),
					},
				},
			}
		},
	},
	Message.MsgType_UserManage_Request: []messageHandler{
		func(msg *Message.Message, sid uint64) *map[uint64]*Message.Message {
			fmt.Printf("MsgType_UserManage_Request (%s)\n", msg.String())
			return &map[uint64]*Message.Message{
				sid: &Message.Message{
					Type:     Message.MsgType_UserManage_Response.Enum(),
					Sequence: proto.Uint32(msg.GetSequence()),
					Response: &Message.Response{
					//Result:           proto.Bool(false),
					//ErrorDescription: []byte("Hello Error"),
					},
				},
			}
		},
	},
}
