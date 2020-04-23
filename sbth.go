package main

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/context"

	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/examples/lib/dev"
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
func Scan(addr string, ctx context.Context) <-chan ble.Advertisement {

	d, err := dev.NewDevice("default")
	if err != nil {
		log.Fatalf("can't new device : %s", err)
	}
	ble.SetDefaultDevice(d)

	// If addr is specified, search for addr instead.
	filter := func(a ble.Advertisement) bool {
		return strings.ToUpper(a.Address().String()) == strings.ToUpper(addr)
	}

	// Scan for specified durantion, or until interrupted by user.
	ch := make(chan ble.Advertisement)
	go func() {
		fn := func(a ble.Advertisement) {
			ch <- a
		}
		err = ble.Scan(ctx ,false, fn, filter)
		if err != nil {
			if err != context.Canceled {
				log.Fatalf("can't connect : %s", err)
			}
		}
	}()
	return ch

}
