package cchelper

import (
	"lib/glog"
	"os"
	"time"
)

type Options struct {
	TCPAddress       string        `flag:"tcp-address"`
	HTTPAddress      string        `flag:"http-address"`
	BroadcastAddress string        `flag:"broadcast-address"`
	TCPReadTimeout   time.Duration `flag:"TCP:net.Conn read timeout"` // 心跳包时间
	// 超时剔除客户端
	//InactiveProducerTimeout time.Duration `flag:"inactive-producer-timeout"`
	//TombstoneLifetime       time.Duration `flag:"tombstone-lifetime"`

	//AcceptTimeout time.Duration
	//ReadTimeout   time.Duration
	//WriteTimeout  time.Duration

	MaxOutstanding int // tcp server 最大吞吐量
}

func NewOptions() *Options {
	hostname, err := os.Hostname()
	if err != nil {
		glog.Fatal(err)
	}

	return &Options{
		TCPAddress:       "0.0.0.0:1250",
		HTTPAddress:      "0.0.0.0:1022",
		BroadcastAddress: hostname,
		TCPReadTimeout:   time.Second * 40, // 默认心跳包30s,客户端如果40s没有发过来心跳包,算作客户端超时离线
		//InactiveProducerTimeout: 300 * time.Second,
		//TombstoneLifetime:       45 * time.Second,
	}
}
