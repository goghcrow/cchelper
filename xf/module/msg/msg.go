package msg

import (
	"lib/Message"
	"lib/proto"
)

func (*MsgModule) GenKeepaliveResp(seq uint32) *Message.Message {
	return &Message.Message{
		Type:     Message.MsgType_Keepalive.Enum(),
		Sequence: proto.Uint32(seq),
	}
}

//func (*MsgModule) GenKeepaliveResp(seq uint32) *Message.Message {

//}
