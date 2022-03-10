package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

//全局chan
var ch = make(chan int64, 20)

func main() {
	//test
	ch <- 100
	ch <- 100
	k, v := <-ch
	fmt.Println(k, v)
	close(ch)
	k, v = <-ch
	fmt.Println(k, v)
	k, v = <-ch
	fmt.Println(k, v)
	os.Exit(1)

	ctx, cancle := context.WithCancel(context.Background())

	go test1(ctx)
	go test2(ctx)

	select {
	case <-time.After(time.Second * 3):
		fmt.Println("cancle")
		cancle()
	}
	for true {
		time.Sleep(time.Second * 100)
		fmt.Println("balbalbalblalbalba")
	}
}

func test1(ctx context.Context) {
	for _ = range time.NewTicker(time.Second * 1).C {
		fmt.Println("test1:", time.Now())
		time.Sleep(time.Second * 1)
		select {
		case <-ctx.Done():
			fmt.Println("test1 done:")
			return
		default:
			fmt.Println("test1 default")
		}
	}
}
func test2(ctx context.Context) {
	for _ = range time.NewTicker(time.Second * 1).C {
		fmt.Println("test2:", <-ch, time.Now())
		time.Sleep(time.Second * 1)
		select {
		case <-ctx.Done():
			fmt.Println("test2 done:")
			return
		default:
			fmt.Println("test2 default")
		}
	}

}
