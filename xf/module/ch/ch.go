package ch

import (
	"lib/link"
	"xf/module"
)

type ChModule struct {
	Global *link.Channel
}

func init() {
	module.Channeler = &ChModule{}
}

func (self *ChModule) xx() {
	//module.Tcp
}
