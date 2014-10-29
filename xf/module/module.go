package module

import (
	"database/sql"
	"lib/Message"
	"lib/link"
	"net"
)

type HttpServer interface {
	Start(listener net.Listener)
	Stop(interface{})
}

type TcpServer interface {
	Start(listener net.Listener)
	Stop(interface{})
	Protocol() link.PacketProtocol
}

type MsgHandler interface {
	// 使用指针的指针来处理无返回值情况 nil
	// A: map[uint64][*Message.Message] key:SentToSid value:SendMsg
	// B: []interface{} 两个一组,基数表示key,偶数表示value 使用时候转型
	// benchmark 表明 性能 A>B
	Handle(*Message.Message, uint64, chan *map[uint64]*Message.Message)
}

type Optioner interface {
	ParseJsonOpt(string) error

	Mysql() string

	Tcp() string
	Http() string
	MaxConn() uint64
	ReadBuffer() int
	ReadTOver() int

	String() string
}

type StatisHandler interface {
	Onlines() int64
	Helps() int
	Publics() int
	HelpSid(sid uint64) uint64
	PublicSid(sid uint64) uint64

	Lock()
	Unlock()
	Online()
	Offline()
	Help(sid uint64)
	Public(sid uint64)
}

type Architecture interface {
	//Sid(string) uint64
	//User(uint64) string
	//UserInfo(string) interface{}

	ArchFetchAll() (map[sql.NullInt64]*ArchRow, error)
	//ArchAppend()
	//ArchUpdate()
	//ArchDel()

	ArchFetchUser() (map[sql.NullString]*EmployeeRow, error)
	//UserAppend()
	//UserUpdate()
	//UserDel()

	//ChanInit()
	//ChanAdd()
	//ChanDel()
}

type Channeler interface {
	//Join(string) // 递归加入
	//Exit(string) // 递归移出来
	//Broadcast(string) //
}

var (
	Opt  Optioner
	Tcp  TcpServer
	Http HttpServer
	Msg  MsgHandler
	Sta  StatisHandler
	Arch Architecture
	Ch   Channeler
)
