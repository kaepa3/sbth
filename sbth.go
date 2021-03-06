package sbth

import (
	"log"
	"strings"
"fmt"
"strconv"

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
	HumidityMask byte = 0x7f
	T1Mask byte = 0x7f
	T2Mask byte = 0x0f

)

func (p *ThermohygroPacket) GetBattery() int {
	return int(p.Packet[2] & BatteryMask)
}

func (p *ThermohygroPacket) GetHumidity() int {
	return int(p.Packet[5] & HumidityMask)
}
func (p *ThermohygroPacket) GetTemperature() float64 {
	t1 := int(p.Packet[4] & T1Mask)
	t2 := int(p.Packet[3] & T2Mask)
	numStr := fmt.Sprintf("%d.%d", t1,t2)
	f, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0 
	}
	return f
}
