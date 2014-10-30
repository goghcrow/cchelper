package msg

import (
	"lib/Message"
	"lib/link"
	"lib/proto"
	"xf/module"
)

var (
	NoHandler = "不能处理..."
)

var mapReq2Res map[Message.MsgType]Message.MsgType = map[Message.MsgType]Message.MsgType{
	Message.MsgType_Login_Request:      Message.MsgType_Login_Response,
	Message.MsgType_Help_Request:       Message.MsgType_Help_Response,
	Message.MsgType_Chat_Request:       Message.MsgType_Chat_Response,
	Message.MsgType_Broadcast_Request:  Message.MsgType_Broadcast_Response,
	Message.MsgType_UserManage_Request: Message.MsgType_UserManage_Response,
}

var errResponse = func(msg *Message.Message, errMsg string, session *link.Session) *[]*module.MsgPack {
	return &[]*module.MsgPack{
		&module.MsgPack{
			Sid: session.Id(),
			Msg: &Message.Message{
				Type:     mapReq2Res[msg.GetType()].Enum(),
				Sequence: proto.Uint32(msg.GetSequence()),
				Response: &Message.Response{
					Result:           proto.Bool(false),
					ErrorDescription: []byte(errMsg),
				},
			},
		},
	}
}
