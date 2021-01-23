package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type newContext chan struct{}

var done = make(chan struct{})

func makeNewContext() context.Context {
	return newContext(done)
}


// Done implements context.Context
func (a newContext) Done() <-chan struct{} {
	return (chan struct{})(a)
}

// Err implements context.Context
func (a newContext) Err() error {
	select {
	case <-(chan struct{})(a):
		return nil
	default:
		return nil
	}
}

// Deadline implements context.Context
func (a newContext) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

// Value implements context.Context
func (a newContext) Value(key interface{}) interface{} {
	return nil
}

func infinite(ctx context.Context) {
	var wg sync.WaitGroup

	go func() {
		for i := 0; i < 2500; i++ {
			fmt.Println("hi")
			time.Sleep(time.Second)
		}
		wg.Done()
	}()
	wg.Add(1)

	wg.Wait()
	fmt.Println("All Done")

}

func restart(close context.CancelFunc) {
	time.Sleep(time.Second*15)
	close()
}

func main() {
	ctx := makeNewContext()
	ctx, cancel := context.WithCancel(ctx)
	go infinite(ctx)
	go restart(cancel)

	<-ctx.Done()

	fmt.Println("Wooo canceled infinite early")


}
