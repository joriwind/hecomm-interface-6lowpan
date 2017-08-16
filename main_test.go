package sixlowpan

import (
	"fmt"
	"net"
	"testing"

	"golang.org/x/net/ipv6"
)

//TestRcvPacket Try to receive a packet via udp-slip connection
/* func TestRcvPacket(t *testing.T) {
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
} */

//Test closing the connection
func TestClose(t *testing.T) {
	cases := []struct {
	}{
		{},
	}
	for _, c := range cases {
		config := Config{
			DebugLevel: DebugAll,
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

//TestContinuous Test the continuous working of connection to udp-slip device
func TestContinuous(t *testing.T) {
	cases := []struct {
	}{
		{},
	}
	for _, c := range cases {
		config := Config{
			DebugLevel: DebugAll,
			PortName:   "/dev/ttyUSB1",
		}
		buf := make([]byte, 200)
		reader, err := Open(config)
		defer reader.Close()

		if err != nil {
			t.Errorf("Did not exit Open properly %v\n", err)
		}
		for {
			n, err := reader.Read(buf)
			if err != nil {
				t.Errorf("Did not exit Read properly %v, %v", err, c)
			}
			fmt.Printf("Packet: %v\n", string(buf[:n]))
		}
	}
}

//TestIPPacket Test the marshalling of IP packet with fixed payload
func TestIPPacket(t *testing.T) {
	h := &ipv6.Header{
		Version:      6,
		TrafficClass: 0,
		FlowLabel:    0,
		PayloadLen:   0,
		NextHeader:   17,
		HopLimit:     255,
	}

	cases := []struct {
		header  *ipv6.Header
		dst     net.IP
		src     net.IP
		payload []byte
	}{
		{header: h, dst: net.ParseIP("aaaa::c30c:0:0:7"), src: net.ParseIP("aaaa::1"), payload: []byte{0x0, 0x0, 0x16, 0x33, 0x0, 0x0A, 0x0, 0x0, 0x48, 0x49}},
	}
	for _, c := range cases {
		c.header.Dst = c.dst
		c.header.Src = c.src
		c.header.PayloadLen = len(c.payload)
		b, err := Marschal(*c.header, c.payload)
		if err != nil {
			t.Errorf("Something went wrong in marshalling IP packet: case %v, result: %x, error: %v\n", c, b, err)
		}
		fmt.Printf("Case: %v, result: %x\n", c, b)

	}
}

//TestUDPPacket Test the Marshalling of UDP packet
func TestUDPPacket(t *testing.T) {
	h := &ipv6.Header{
		Version:      6,
		TrafficClass: 0,
		FlowLabel:    0,
		PayloadLen:   0x001e,
		NextHeader:   17,
		HopLimit:     63,
		Src:          net.ParseIP("aaaa::c30c:0:0:7"),
		Dst:          net.ParseIP("aaaa::1"),
	}
	uHdr := &UDPHeader{
		DstPort: 0xb25f,
		SrcPort: 0x1633,
		Length:  0x001e,
		Chksum:  0,
		Payload: []byte{0x60, 0x45, 0x22, 0xa6, 0x41, 0x0c, 0x80, 0xb1, 0x02, 0xff, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21},
	}

	cases := []struct {
		header *ipv6.Header
		udpHdr *UDPHeader
		result []byte
	}{
		{
			header: h,
			udpHdr: uHdr,
			result: []byte{0x16, 0x33, 0xb2, 0x5f, 0x00, 0x1e, 0x85, 0x1e, 0x60, 0x45, 0x22, 0xa6,
				0x41, 0x0c, 0x80, 0xb1, 0x02, 0xff, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21}},
	}
	for _, c := range cases {
		err := c.udpHdr.CalcChecksum(c.header)
		if err != nil {
			t.Errorf("Something went wrong in calculating checksum UDP packet: case %v, error: %v\n", c, err)
		}
		b, err := c.udpHdr.Marschal()
		if err != nil {
			t.Errorf("Something went wrong in marshalling UDP packet: case %v, result: %x, error: %v\n", c, b, err)
		}
		if c.result != nil {
			if len(c.result) != len(b) {
				t.Errorf("Did not correspond to expected result, size error!\n")
			}
			for i, x := range b {
				if c.result[i] != x {
					t.Errorf("Did not correspond to expected result, content error!\n")
				}
			}
		}
		fmt.Printf("Succesful:: Case: %v, result: %x\n", c, b)

	}
}

//TestIpUDPPacket Test the full range of UDP to IP to bytes working
func TestIPUDPPacket(t *testing.T) {
	h := &ipv6.Header{
		Version:      6,
		TrafficClass: 0,
		FlowLabel:    0,
		PayloadLen:   0x001e,
		NextHeader:   17,
		HopLimit:     63,
		Src:          net.ParseIP("aaaa::c30c:0:0:7"),
		Dst:          net.ParseIP("aaaa::1"),
	}
	uHdr := &UDPHeader{
		DstPort: 0xb25f,
		SrcPort: 0x1633,
		Length:  0x001e,
		Chksum:  0,
		Payload: []byte{0x60, 0x45, 0x22, 0xa6, 0x41, 0x0c, 0x80, 0xb1, 0x02, 0xff, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21},
	}

	cases := []struct {
		header *ipv6.Header
		udpHdr *UDPHeader
		result []byte
	}{
		{
			header: h,
			udpHdr: uHdr,
			result: []byte{
				0x60, 0x00, 0x00, 0x00, 0x00, 0x1e, 0x11, 0x3f,
				0xaa, 0xaa, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc3, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x07,
				0xaa, 0xaa, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
				0x16, 0x33, 0xb2, 0x5f, 0x00, 0x1e, 0x85, 0x1e,
				0x60, 0x45, 0x22, 0xa6, 0x41, 0x0c, 0x80, 0xb1, 0x02, 0xff, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21}},
	}
	for _, c := range cases {
		err := c.udpHdr.CalcChecksum(c.header)
		if err != nil {
			t.Errorf("Something went wrong in calculating checksum UDP packet: case %v, error: %v\n", c, err)
		}
		b, err := c.udpHdr.Marschal()
		if err != nil {
			t.Errorf("Something went wrong in marshalling UDP packet: case %v, result: %x, error: %v\n", c, b, err)
		}

		fullIP, err := Marschal(*c.header, b)
		if err != nil {
			t.Errorf("Something went wrong in marshalling IP packet: case %v, result: %x, error: %v\n", c, b, err)
		}

		if c.result != nil {
			if len(c.result) != len(fullIP) {
				t.Errorf("Did not correspond to expected result, size error!\n")
			}
			for i, x := range fullIP {
				if c.result[i] != x {
					t.Errorf("Compare result:: content error: index: %v, expected: %x, got: %x!\n", i, c.result[i], x)
				}
			}
		}
		fmt.Printf("Succesful:: Case: %v, result: %x\n", c, fullIP)

	}
}
