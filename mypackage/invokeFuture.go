package mypackage

import (
	"fmt"
	"sync"
	"testing"
)

type InvokeFuture interface {
	finish() int
}
type DefaultInvokeFuture struct {
	InvokeId  int
	MyChannel chan int
}

var FutureMap sync.Map

// 工厂模式
func (t *DefaultInvokeFuture) finish() int {
	x := <-t.MyChannel
	return x
}

func TestFuture(t *testing.T) {

	//suppose send msg in main thread.

	ch := make(chan int)

	invokeId := 10

	x := DefaultInvokeFuture{invokeId, ch}

	//put x into futureMap
	FutureMap.Store(invokeId, x)

	go func() {
		//receive the msg ,put this on ch
		fmt.Println("receive msg")
		value, ok := FutureMap.LoadAndDelete(invokeId)

		if ok {
			future := value.(DefaultInvokeFuture)
			future.MyChannel <- 55
		}
	}()

	// sync receive the msg
	result := <-ch

	fmt.Println(result)

}
