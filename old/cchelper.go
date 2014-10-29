package cchelper

import (
	"lib/glog"
	"net"
	"os"
	"sync"
)

// 配置+tcp服务+消息服务+http服务+waitgroup
type CCHelper struct {
	opts         *Options
	tcpAddr      *net.TCPAddr
	tcpListener  net.Listener
	tcpMsgChan   chan *chanMsgReq // 发送消息至消息处理循环
	httpAddr     *net.TCPAddr
	httpListener net.Listener
	wg           sync.WaitGroup
}

type CCHCtx struct{ ctx *CCHelper } // CCHelper Context

//func init() {
//flag.Parse() // 从这里删除,放在main包main方法中
//}

func New(opts *Options) *CCHelper {
	//defer glog.Flush() // 从这里删除,放在main包main方法中

	cch := &CCHelper{opts: opts}
	tcpAddr, err := net.ResolveTCPAddr("tcp", opts.TCPAddress)
	fatalErrCheck(err)
	cch.tcpAddr = tcpAddr

	httpAddr, err := net.ResolveTCPAddr("tcp", opts.HTTPAddress)
	fatalErrCheck(err)
	cch.httpAddr = httpAddr

	return cch
}

func (cch *CCHelper) Start() {
	glog.Info("=== >CCHelper Service< ===")

	defer func() {
		if err := recover(); err != nil {
			glog.Fatalf("CChelper panic(%s)", err)
		}
	}()

	// 数据报告
	go func() {
		cch.wg.Add(1)
		glog.Info("StatisServer start")
		StatisServer()
		glog.Info("StatisServer end")
		cch.wg.Done()
	}()

	// 消息处理
	cch.tcpMsgChan = make(chan *chanMsgReq)

	go func() {
		cch.wg.Add(1)
		glog.Info("TCPMessageLoop start")
		TCPMessageLoop(cch.tcpMsgChan, &messageServer{ctx: cch})
		glog.Info("TCPMessageLoop end")
		cch.wg.Done()
	}()

	// tcp服务器
	tcpListener, err := net.Listen("tcp", cch.tcpAddr.String())
	fatalErrCheck(err)
	cch.tcpListener = tcpListener

	go func() {
		cch.wg.Add(1)
		glog.Info("TCPServer start")
		TCPServer(tcpListener, &tcpServer{ctx: cch})
		glog.Info("TCPServer end")
		cch.wg.Done()
	}()

	// http服务器
	httpListener, err := net.Listen("tcp", cch.httpAddr.String())
	fatalErrCheck(err)
	cch.httpListener = httpListener

	go func() {
		cch.wg.Add(1)
		glog.Info("HTTPServer start")
		HTTPServer(httpListener, &httpServer{ctx: cch})
		glog.Info("HTTPServer end")
		cch.wg.Done()
	}()
}

func (cch *CCHelper) Exit() {
	if cch.tcpListener != nil {
		cch.tcpListener.Close()
	}
	if cch.httpListener != nil {
		cch.httpListener.Close()
	}

	close(cch.tcpMsgChan) // 未判断是否已经关闭

	cch.wg.Wait()
}

func fatalErrCheck(err error) {
	if err != nil {
		glog.Fatal(err)
		glog.Flush()
		os.Exit(1)
	}
}
