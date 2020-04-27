package sbth

import (
	"log"
	"strings"

	"golang.org/x/net/context"

	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/examples/lib/dev"
)

func Scan(addr string, ctx context.Context) <-chan ThermohygroPacket {

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
	ch := make(chan ThermohygroPacket)
	go func() {
		fn := func(a ble.Advertisement) {
			b := a.ServiceData()
			if 1 == len(b) {
				pack := ThermohygroPacket{[]byte(b[0].UUID), b[0].Data}
				ch <- pack
			}

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

type ThermohygroData []byte

type ThermohygroPacket struct {
	Uuid   []byte
	Packet ThermohygroData
}

const (
	BatteryMask byte = 0x7f
)

func (p *ThermohygroPacket) GetBattery() int {
	return int(p.Packet[2] & BatteryMask)
}
