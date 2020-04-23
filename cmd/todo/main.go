package main

import (
	"context"
	"fmt"
)

func main() {
	addr := "DB:FA:C9:3C:48:A2"
	ctx, _ := context.WithCancel(context.Background())
	ch := Scan(addr, ctx)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case p := <-ch:
				fmt.Println("rcv")
				fmt.Println(p)
				ctx.Done()
				close(done)
				break
			case <-ctx.Done():
				close(done)
				break
			}
		}
	}()
	<-done
}
