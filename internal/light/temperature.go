package light

import (
	"errors"
	"net"
	"strconv"
)

// changeLightTemperature sends a message to the light bulb to change it's temperature. The temperature parameter
// is in kelvin (e.g. 1000, 4500, 6000, etc). lightIP is the ip address of the bulb.
func ChangeLightTemperature(conn net.PacketConn, lightIP net.Addr, temperature int) error {
	if temperature < 0 {
		return errors.New("cannot set negative temperature value")
	}
	message := `{"method": "setPilot", "params": {"temp":` + strconv.Itoa(temperature) + `}}`
	_, err := conn.WriteTo([]byte(message), lightIP)
	return err
}
