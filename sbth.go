package sbth

import (
	"log"
	"strings"

	"golang.org/x/net/context"

	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/examples/lib/dev"
)

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
		err = ble.Scan(ctx, false, fn, filter)
		if err != nil {
			if err != context.Canceled {
				log.Fatalf("can't connect : %s", err)
			}
		}
	}()
	return ch

}
