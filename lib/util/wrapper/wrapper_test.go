package wrapper

import (
	"fmt"
	"testing"
)

func TestWaitGroup(t *testing.T) {
	wg := NewWaitGroup()
	wg.AsynRun(func() {
		fmt.Println("AsynRun")
	})
	wg.Wait()
}
