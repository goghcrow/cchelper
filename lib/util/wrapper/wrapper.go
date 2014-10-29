package wrapper

import "sync"

type WaitGroup struct {
	sync.WaitGroup
}

func NewWaitGroup() *WaitGroup { return &WaitGroup{} }

func (wg *WaitGroup) AsynRun(fun func()) {
	wg.Add(1)
	go func() {
		fun()
		wg.Done()
	}()
}

//func (wg *WaitGroup) Wait() {
//	wg.Wait()
//}
