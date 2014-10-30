package statis

import (
	"lib/link"
	"sync"
	"sync/atomic"
	"xf/module"
)

type StatisModule struct {
	rw          sync.RWMutex
	OnlineCount int64 // uint64 can't atomic.AddUint64(&self.OnlineCount, -1)
	HelpCount   map[uint64]*uint64
	PublicCount map[uint64]*uint64
}

func init() {
	module.Sta = &StatisModule{
		HelpCount:   make(map[uint64]*uint64),
		PublicCount: make(map[uint64]*uint64),
	}
}
func (self *StatisModule) Onlines() int64 { return self.OnlineCount }
func (self *StatisModule) Helps() int     { return len(self.HelpCount) }
func (self *StatisModule) Publics() int   { return len(self.PublicCount) }
func (self *StatisModule) HelpSid(sid uint64) uint64 {
	if p, ok := self.HelpCount[sid]; ok {
		return *p
	} else {
		return 0
	}
}
func (self *StatisModule) PublicSid(sid uint64) uint64 {
	if p, ok := self.PublicCount[sid]; ok {
		return *p
	} else {
		return 0
	}
}

func (self *StatisModule) Lock()    { self.rw.Lock() }
func (self *StatisModule) Unlock()  { self.rw.Unlock() }
func (self *StatisModule) Online()  { atomic.AddInt64(&self.OnlineCount, 1) }
func (self *StatisModule) Offline() { atomic.AddInt64(&self.OnlineCount, -1) }
func (self *StatisModule) Help(session *link.Session) {
	sid := session.Id()
	self.rw.Lock()
	var pInt64 *uint64
	pInt64, ok := self.HelpCount[sid]
	if !ok {
		var t uint64
		pInt64 = &t
		self.HelpCount[sid] = &t
	}
	atomic.AddUint64(pInt64, 1)
	self.rw.Unlock()
}
func (self *StatisModule) Public(session *link.Session) {
	sid := session.Id()
	self.rw.Lock()
	var pInt64 *uint64
	pInt64, ok := self.PublicCount[sid]
	if !ok {
		var t uint64
		pInt64 = &t
		self.PublicCount[sid] = &t
	}
	atomic.AddUint64(pInt64, 1)
	self.rw.Unlock()
}
