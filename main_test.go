package sixlowpan

import (
	"fmt"
	"testing"
)

func TestRcvPacket(t *testing.T) {
	cases := []struct {
	}{
		{},
	}
	for _, c := range cases {
		config := Config{
			DebugLevel: DebugPacket,
			PortName:   "/dev/ttyUSB1",
		}
		buf := make([]byte, 200)
		reader, err := Open(config)
		defer reader.Close()

		if err != nil {
			t.Errorf("Did not exit Open properly %v\n", err)
		}
		n, err := reader.Read(buf)
		fmt.Printf("Packet: %v", string(buf[:n]))
		if err != nil {
			t.Errorf("Did not exit Read properly %v, %v", err, c)
		}
	}
}

func TestClose(t *testing.T) {
	cases := []struct {
	}{
		{},
	}
	for _, c := range cases {
		config := Config{
			DebugLevel: DebugPacket,
			PortName:   "/dev/ttyUSB1",
		}
		buf := make([]byte, 20)
		reader, err := Open(config)
		if err != nil {
			t.Errorf("Did not exit Open properly %v\n", err)
		}

		reader.Close()

		_, err = reader.Read(buf)
		if err == nil {
			t.Errorf("Did not Close properly %v, %v", err, c)
		}
	}
}
