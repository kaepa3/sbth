package main

import (
	"context"
	"fmt"
	"time"

	"github.com/kaepa3/sbth"
)

func main() {
	addr := "DB:FA:C9:3C:48:A2"
	ctx, _:= context.WithCancel(context.Background())
	timer := time.NewTimer(time.Second * 3)
	ch := sbth.Scan(addr, ctx)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case p := <-ch:
				fmt.Println("rcv")
				fmt.Println(p.GetBattery())
				fmt.Println(p.GetHumidity())
				fmt.Println(p.GetTemperature())
				break
			case <-ctx.Done():
			case <-timer.C:
				close(done)
				break
			}
		}
	}()
	<-done
}
