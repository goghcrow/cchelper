// X添加channel支持
// X添加统计
// 目前的消息接口不适合广播....
// SendQueue
// BUG: 断线重连之后,读超时失败
// 广播时候 seq 填写最大值或者最小值,或者在客户端规定广播seq不处理

// todo: 手动做一下内联...将tcp模块三个文件合并
// todo: 使用server.Send来处理全部发送部分,舍弃session.Send()
// todo: msgmodule与tcpmodule合并//用新版合并
// todo: 是否需要统一处理发送
// todo: packet head 写入配置文件
package main

import (
	"bufio"
	"fmt"
	"lib/glog"
	"lib/util/wrapper"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	"xf/module"
)

import (
	_ "xf/module/arch"
	_ "xf/module/http"
	_ "xf/module/msg"
	_ "xf/module/opt"
	_ "xf/module/statis"
	_ "xf/module/tcp"
)

var wg = wrapper.NewWaitGroup()

func main() {
	defer glog.Flush()

	glog.Info("xf server start")
	glog.Info(module.Opt.String())

	//arches, err := module.Arch.ArchFetchAll()
	//fatalErrCheck(err)
	//for id, row := range arches {
	//	fmt.Printf("[%d] {id:%d|pid:%d|depth:%d|order:%d|path:%s|name:%s}\n",
	//		id.Int64,
	//		row.Id.Int64,
	//		row.Pid.Int64,
	//		row.Depth.Int64,
	//		row.Order.Int64,
	//		row.Path.String,
	//		row.Name.String,
	//	)
	//}

	//users, err := module.Arch.ArchFetchUser()
	//fatalErrCheck(err)
	//for user, row := range users {
	//	fmt.Printf("[%s] {id:%d|archid:%d|erp:%s|user:%s|name:%s|identity:%s}\n",
	//		user.String,
	//		row.Id.Int64,
	//		row.ArchId.Int64,
	//		row.Erp.String,
	//		row.User.String,
	//		row.Name.String,
	//		row.Identity.String,
	//	)
	//}

	// start listen tcp
	wg.AsynRun(func() {
		defer module.Tcp.Stop(nil)
		tcpListener, err := net.Listen("tcp", module.Opt.Tcp())
		fatalErrCheck(err)
		module.Tcp.Start(tcpListener)
	})

	// start listen http // 别忘了实现Stop方法...
	//wg.AsynRun(func() {
	//	defer module.Http.Stop(nil)
	//	httpListener, err := net.Listen("tcp", module.Opt.Http())
	//	fatalErrCheck(err)
	//	module.Http.Start(httpListener)
	//})

	exit := false
	reader := bufio.NewReader(os.Stdin)
	for !exit {
		data, _, _ := reader.ReadLine()
		cmds := strings.SplitN(string(data), " ", 2)
		switch cmds[0] {
		case "x": // exit
			Exit()
			exit = true
		case "r": //report
			report()
		case "h":
			if len(cmds) < 2 {
				continue
			}
			sid, err := strconv.ParseUint(cmds[1], 10, 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			help(sid)
		}
	}

}

func Exit() {
	module.Tcp.Stop("Manual Exit")
	module.Http.Stop("Manual Exit")
	wg.Wait()
}

func fatalErrCheck(err error) {
	if err != nil {
		glog.Fatal(err)
		glog.Flush()
		os.Exit(1)
	}
}

func help(sid uint64) {
	fmt.Printf(`[%s]
++++++++++++++++++++++++++++++++++++++++
Sid : [%d]
	Help   : [%d]
	Public : [%d]
++++++++++++++++++++++++++++++++++++++++
`, time.Now().String(), sid, module.Sta.HelpSid(sid), module.Sta.PublicSid(sid))
}

func report() {
	fmt.Printf(
		`[%s]
++++++++++++++++++++++++++++++++++++++++
Online : [%d]
----------------------------------------
Help
	HelpLen   : [%d]
	PublicLen : [%d]
----------------------------------------
Clients
	UuidLen : [%d]
	NameLen : [%d]
++++++++++++++++++++++++++++++++++++++++
`,
		time.Now().String(),
		module.Sta.Onlines(),
		module.Sta.Helps(),
		module.Sta.Publics(),
		0,
		0,
	)
}
