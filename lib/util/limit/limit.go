package limit

type null struct{}
type MaxBlockCtrl struct {
	ch  chan null
	max uint64
}

func NewMaxBlockCtrl(max uint64) *MaxBlockCtrl {
	mbc := &MaxBlockCtrl{ch: make(chan null, max), max: max}
	var i uint64 = 0
	for ; i < max; i++ {
		mbc.ch <- null{}
	}
	return mbc
}

func (self *MaxBlockCtrl) Pop()              { <-self.ch }
func (self *MaxBlockCtrl) Push()             { self.ch <- null{} }
func (self *MaxBlockCtrl) Max() uint64       { return self.max }
func (self *MaxBlockCtrl) SetMax(max uint64) { self.ch = make(chan null, max) }
