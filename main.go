package main

import (
	"context"
	"fmt"
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

func infTask(ctx context.Context, myNum int) {
	for {
		select {
		case <-ctx.Done(): //if Done channel closed exit the function
			return
		default:
			//keep going
		}

		time.Sleep(time.Second)
		fmt.Println(myNum)
	}
}

// startTasks spawns 20 goroutines each with it's own infinite task
// We use the context.CancelFunc() function in restart() to stop the work of all associated goroutines early
func startTasks(ctx context.Context) {
	for i := 0; i < 20; i++ {
		go infTask(ctx, i)
	}
}

func restart(close context.CancelFunc) {
	time.Sleep(time.Second * 10)
	close()
}

func main() {
	ctx := makeNewContext()
	ctx, cancel := context.WithCancel(ctx)
	go startTasks(ctx)
	go restart(cancel)

	//Blocks unitl this context's Done channel is closed
	//(calling context.CancelFunc like in restart closes the Done channel)
	<-ctx.Done()

	fmt.Println("Wooo canceled infinite early")

	go func() {
		for {
			time.Sleep(time.Second)
		}
	}()

	select {} //block main infinetly to make sure nothing more is printed and all the
	//infTask() go routines have returned
}
