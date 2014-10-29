package limit

import (
	"lib/dbgutil"
	"testing"
	"time"
)

func TestMaxBlockCtrl(t *testing.T) {
	sem := NewMaxBlockCtrl(10)

	pCount := 0
	go func() {
		for {
			sem.Pop()
			pCount++
			dbgutil.FormatDisplay("pCount", pCount)
		}
	}()

	go func() {
		ch := time.Tick(time.Second)
		for {
			<-ch
			sem.Push()
		}
	}()

	select {}
}
