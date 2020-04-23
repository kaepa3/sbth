package sbth

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/context"

	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/examples/lib/dev"
)

func Scan(addr string, ctx context.Context) {

	d, err := dev.NewDevice("default")
	if err != nil {
		log.Fatalf("can't new device : %s", err)
	}
	ble.SetDefaultDevice(d)

	// If addr is specified, search for addr instead.
	filter := func(a ble.Advertisement) bool {
		flg := strings.ToUpper(a.Address().String()) == strings.ToUpper(addr)
		if flg {
			fmt.Println(a.ServiceData())
		}
		return flg
	}

	ctx2, cancel := context.WithCancel(ctx)
	go func() {
		select {
		case <-ctx.Done():
			cancel()
		case <-ctx2.Done():
		}
	}()

	// Scan for specified durantion, or until interrupted by user.
	ctxBle := ble.WithSigHandler(ctx2, cancel)
	cln, err := ble.Connect(ctxBle, filter)
	if err != nil {
		log.Fatalf("can't connect : %s", err)
	}

	// Make sure we had the chance to print out the message.
	done := make(chan struct{})

	// Normally, the connection is disconnected by us after our exploration.
	// However, it can be asynchronously disconnected by the remote peripheral.
	// So we wait(detect) the disconnection in the go routine.
	go func() {
		<-cln.Disconnected()
		fmt.Printf("[ %s ] is disconnected \n", cln.Address())
		close(done)
	}()

	fmt.Printf("Discovering profile...\n")
	_, err = cln.DiscoverProfile(true)
	if err != nil {
		log.Fatalf("can't discover profile: %s", err)
	}

	// Start the exploration.

	// Disconnect the connection. (On OS X, this might take a while.)
	fmt.Printf("Disconnecting [ %s ]... (this might take up to few seconds on OS X)\n", cln.Address())
	cln.CancelConnection()

	<-done
}
