// two ways conf [flag | conf.json]
package opt

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"xf/module"
)

type OptionModule struct {
	json             bool   `override flag`
	path             string `conf.json path`
	TcpAddr          string `TcpAddr`
	HttpAddr         string `HttpAddr`
	TcpMaxConnection uint64 `TcpMaxConnection`
	ReadBufferSize   int    `ReadBufferSize`
	ReadOverTime     int    `ReadOverTime`
	MysqlDSN         string `Mysql Data Source Name`
}

func init() {
	opt := &OptionModule{}

	flag.BoolVar(&opt.json, "json", false, "use conf.json")
	flag.StringVar(&opt.path, "path", "conf.json", "path of conf.json")
	flag.StringVar(&opt.TcpAddr, "tcpaddr", "127.0.0.1:1022", "tcp server listening addaddress")
	flag.StringVar(&opt.HttpAddr, "httpaddr", "127.0.0.1:80", "http server listening addaddress")
	flag.Uint64Var(&opt.TcpMaxConnection, "tcpmaxconnection", 50000, "max count of tcp keep-alive connection")
	flag.IntVar(&opt.ReadBufferSize, "readbuffersize", 1024, "tcp read buffer size")
	flag.IntVar(&opt.ReadOverTime, "readovertime", 40, "tcp read time over(second)")
	flag.StringVar(&opt.MysqlDSN, "mysql", "root@/", "Mysql data source name")
	flag.Parse() // the same as glog parse()

	if opt.json {
		err := opt.ParseJsonOpt(opt.path)
		if err != nil {
			fmt.Printf("json option file parse err [%s]", err.Error())
			os.Exit(-1)
		}
	}

	module.Opt = opt
}

func (opt *OptionModule) ParseJsonOpt(path string) error {
	bf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bf, &opt)
	if err != nil {
		return err
	}
	return nil
}

func (opt *OptionModule) Tcp() string     { return opt.TcpAddr }
func (opt *OptionModule) Http() string    { return opt.HttpAddr }
func (opt *OptionModule) MaxConn() uint64 { return opt.TcpMaxConnection }
func (opt *OptionModule) ReadBuffer() int { return opt.ReadBufferSize }
func (opt *OptionModule) ReadTOver() int  { return opt.ReadOverTime }
func (opt *OptionModule) Mysql() string   { return opt.MysqlDSN }

func (opt *OptionModule) String() string {
	var typ string
	if opt.json {
		typ = "\n[conf.json "
	} else {
		typ = "\n[flag "
	}
	return typ + fmt.Sprintf(`Options]
	TcpAddr          %s
	HttpAddr         %s
	TcpMaxConnection %d
	ReadBufferSize   %d
	ReadTimeOver     %d
	MysqlDSN         %s
`, opt.TcpAddr, opt.HttpAddr, opt.TcpMaxConnection, opt.ReadBufferSize, opt.ReadOverTime, opt.MysqlDSN)
}
