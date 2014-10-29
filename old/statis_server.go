package cchelper

import (
	"sync"
	"sync/atomic"
)

// region global
var uuid int64

func Uuid() int64 {
	return atomic.AddInt64(&uuid, 1)
}

// 从架构里拆出来 可以提高range map 性能...
type clients struct {
	rwm    sync.RWMutex
	Uuid2c map[int64]*client // Uuid2c 与 Name2c 指向同一个client
	Name2c map[string]*client
}

type client struct {
	//online       bool
	username     string
	chFromMarMsg chan []byte
}

var Clients clients = clients{Uuid2c: make(map[int64]*client), Name2c: make(map[string]*client)}

// endregion global

var (
	OnlineCount int64 = 0
	Help_             = map[string]*int64{}
	Help_Public       = map[string]*int64{}
)

var (
	chOnline chan int64 = make(chan int64) // online <- 1 offline <- -1
	chHelp   chan int64 = make(chan int64) // uuid int64
	chPublic chan int64 = make(chan int64) // uuid int64
)

var StatisHandler statisHandler = map[chan int64][]func(int64){
	chOnline: []func(int64){
		func(onoff int64) {
			atomic.AddInt64(&OnlineCount, onoff)
		},
	},
	chHelp: []func(int64){
		func(uuid int64) {
			if c, ok := Clients.Uuid2c[uuid]; ok {
				var pInt64 *int64
				pInt64, ok := Help_[c.username]
				if !ok {
					t := int64(0)
					pInt64 = &t
					Help_[c.username] = &t
				}
				atomic.AddInt64(pInt64, 1)
			}
		},
	},
	chPublic: []func(int64){
		func(uuid int64) {
			if c, ok := Clients.Uuid2c[uuid]; ok {
				var pInt64 *int64
				pInt64, ok := Help_Public[c.username]
				if !ok {
					t := int64(0)
					pInt64 = &t
					Help_Public[c.username] = &t
				}
				atomic.AddInt64(pInt64, 1)
			}
		},
	},
}

// region statisServer
type statisHandler map[chan int64][]func(int64) //uuid

func (s *statisHandler) Register(ch chan int64, fn func(int64)) {
	if _, ok := StatisHandler[ch]; ok {
		StatisHandler[ch] = append(StatisHandler[ch], fn)
	} else {
		StatisHandler[ch] = []func(int64){fn}
	}
}

func (s statisHandler) Range(ch chan int64, uuid int64) {
	for _, handler := range s[ch] {
		handler(uuid)
	}
}

func StatisServer() {
	defer func() {
		close(chOnline)
		close(chHelp)
		close(chPublic)
	}()

	for {
		select {
		case onoff := <-chOnline:
			StatisHandler.Range(chOnline, onoff)
		case uuid := <-chHelp:
			StatisHandler.Range(chHelp, uuid)
		case uuid := <-chPublic:
			StatisHandler.Range(chPublic, uuid)
		}
	}
}

// endregion statisServer
