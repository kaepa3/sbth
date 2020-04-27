package sbth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetBatteryTest(t *testing.T) {
	packet := ThermohygroPacket{[]byte{0x00}, []byte{0, 0, 0x82}}
	assert.Equal(t, 2, packet.GetBattery())

}

func GetHumidityTest(t *testing.T) {
	packet := ThermohygroPacket{[]byte{0x00}, []byte{0, 0, 0, 0, 0, 100}}
	assert.Equal(t, 100, packet.GetHumidity())
}
func GetTemperature(t *testing.T) {
	packet := ThermohygroPacket{[]byte{0x00}, []byte{0, 0, 0, 4, 30, 100}}
	assert.Equal(t, 30.4, packet.GetHumidity())
}
