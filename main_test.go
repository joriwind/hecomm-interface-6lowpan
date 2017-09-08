package sixlowpan

import (
	"bytes"
	"fmt"
	"net"
	"testing"
	"time"

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

//TestSndPacket Try to receive a packet via udp-slip connection
/* func TestSndPacket(t *testing.T) {
	payload := []byte{0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21}
	ipHdr := &ipv6.Header{
		Version:      6,
		TrafficClass: 0,
		FlowLabel:    0,
		PayloadLen:   len(payload) + UdpHeaderLen,
		NextHeader:   17,
		HopLimit:     63,
		Src:          net.ParseIP("aaaa::c30c:0:0:7"),
		Dst:          net.ParseIP("aaaa::1"),
	}

	udpHdr := &UDPHeader{
		DstPort: 0xb25f,
		SrcPort: 0x1633,
		Length:  uint16(len(payload) + UdpHeaderLen),
		Chksum:  0,
		Payload: payload,
	}

	err := udpHdr.CalcChecksum(ipHdr)
	if err != nil {
		t.Errorf("Something went wrong in calculating checksum UDP packet: ipHdr %v, udpHdr: %v, error: %v\n", ipHdr, udpHdr, err)
	}
	udpb, err := udpHdr.Marschal()
	if err != nil {
		t.Errorf("Something went wrong in marshalling UDP packet: ipHdr %v, udpHdr: %v, error: %v\n", ipHdr, udpHdr, err)
	}

	fullIP, err := Marschal(*ipHdr, udpb)
	if err != nil {
		t.Errorf("Something went wrong in marshalling IP packet: ipHdr %v, udpbytes: %v, error: %v\n", ipHdr, udpb, err)
	}

	cases := []struct {
		input []byte
	}{
		{input: fullIP},
	}
	for _, c := range cases {
		config := Config{
			DebugLevel: DebugAll,
			PortName:   "/dev/ttyUSB1",
		}
		writer, err := Open(config)
		defer writer.Close()

		if err != nil {
			t.Errorf("Did not exit Open properly %v\n", err)
		}
		n, err := writer.Write(c.input)
		if err != nil {
			t.Errorf("Did not exit Write properly %v, %v\n", err, c)
		}
		fmt.Printf("Written: n: %v, input: %x\n", n, c.input)
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
	//Construct package to write
	payload := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x0a} //, 0x54, 0x45, 0x53, 0x54
	ipHdr := ipv6.Header{
		Version:      6,
		TrafficClass: 0,
		FlowLabel:    5,
		PayloadLen:   len(payload) + UdpHeaderLen,
		NextHeader:   17,
		HopLimit:     63,
		Src:          net.ParseIP("aaaa::c30c:0:0:7"),
		Dst:          net.ParseIP("aaaa::1"),
	}

	udpHdr := &UDPHeader{
		DstPort: 0xb25f,
		SrcPort: 0x1633,
		Length:  uint16(len(payload) + UdpHeaderLen),
		Chksum:  0,
		Payload: payload,
	}

	err := udpHdr.CalcChecksum(ipHdr)
	if err != nil {
		t.Errorf("Something went wrong in calculating checksum UDP packet: ipHdr %v, udpHdr: %v, error: %v\n", ipHdr, udpHdr, err)
	}
	udpb, err := udpHdr.Marschal()
	if err != nil {
		t.Errorf("Something went wrong in marshalling UDP packet: ipHdr %v, udpHdr: %v, error: %v\n", ipHdr, udpHdr, err)
	}

	fullIP, err := Marschal(ipHdr, udpb)
	if err != nil {
		t.Errorf("Something went wrong in marshalling IP packet: ipHdr %v, udpbytes: %v, error: %v\n", ipHdr, udpb, err)
	}

	//Startup interface
	config := Config{
		DebugLevel: DebugAll,
		PortName:   "/dev/ttyUSB1",
	}
	buf := make([]byte, 200)
	readwriteclose, err := Open(config)
	defer readwriteclose.Close()

	if err != nil {
		t.Errorf("Did not exit Open properly %v\n", err)
	}

	command := "W"
	duration := 2 * time.Second
	go func() {
		buf := make([]byte, 120)
		for {
			readwriteclose.Read(buf)
		}
	}()
	//Looping
	for {

		time.Sleep(duration)

		switch command {
		case "R":
			n, err := readwriteclose.Read(buf)
			if err != nil {
				t.Errorf("Did not exit Read properly %v\n", err)
			}
			fmt.Printf("Packet: %x\n", string(buf[:n]))

		case "W":
			n, err := readwriteclose.Write(fullIP)
			if err != nil {
				t.Errorf("Did not exit Write properly %v, input: %x\n", err, fullIP)
			}
			fmt.Printf("Packet: %x\n", string(fullIP[:n]))
		case "C":
			readwriteclose.Close()
			return
		default:
			//fmt.Printf("-> Did not understand command: %v\n", command)
		}

	}
}

//TestIPPacket Test the marshalling of IP packet with fixed payload
func TestIPPacket(t *testing.T) {
	h := ipv6.Header{
		Version:      6,
		TrafficClass: 0,
		FlowLabel:    0,
		PayloadLen:   0,
		NextHeader:   17,
		HopLimit:     255,
		Dst:          net.ParseIP("aaaa::c30c:0:0:7"),
		Src:          net.ParseIP("aaaa::1"),
	}

	cases := []struct {
		ipheader  ipv6.Header
		udppacket UDPHeader
	}{
		{ipheader: h,
			udppacket: UDPHeader{SrcPort: 1520, DstPort: 5683, Length: 0, Chksum: 0, Payload: []byte{0x48, 0x49, 0x0a}}},
	}
	for _, c := range cases {
		c.udppacket.Length = uint16(len(c.udppacket.Payload)) + UdpHeaderLen
		c.ipheader.PayloadLen = int(c.udppacket.Length)
		c.udppacket.CalcChecksum(c.ipheader)

		ipPayload, err := c.udppacket.Marschal()
		if err != nil {
			t.Errorf("Unable to marshal UDP packet: %v\n", err)
			return
		}
		b, err := Marschal(c.ipheader, ipPayload)

		ipheader, err := ipv6.ParseHeader(b[:ipv6.HeaderLen])
		if err != nil {
			fmt.Printf("Error in processing ipv6 header\n")
			t.Errorf("IPV6 parse failed: %v\n", err)
			return

		}
		fmt.Printf("Parsed IP header: %v\n", ipheader)
		udpheader, err := UnmarshalUDP(b[ipv6.HeaderLen:])
		if err != nil {
			t.Errorf("UDP parse failed: %v\n", err)
			return
		}
		fmt.Printf("Parsed UDP header: %v\n", udpheader)
		if !bytes.Equal(udpheader.Payload, c.udppacket.Payload) {
			t.Errorf("Test case payload does not match: case: %v, generated: %v\n", c.udppacket.Payload, udpheader.Payload)
			return
		}
		//Checking checksum
		err = udpheader.CalcChecksum(*ipheader)
		if err != nil {
			t.Errorf("Could not calculate checksum: udp: %v, ip: %v", udpheader, ipheader)
		}
		if udpheader.Chksum != c.udppacket.Chksum {
			t.Errorf("Test case checksum does not match: expected: %v, generated: %v\n", c.udppacket.Chksum, udpheader.Chksum)
			return
		}
		if err != nil {
			t.Errorf("Something went wrong in marshalling IP packet: case %v, result: %x, error: %v\n", c, b, err)
			return
		}
		//fmt.Printf("Case: %v, result: %x\n", c, b)

	}
}

//TestUDPPacket Test the Marshalling of UDP packet
func TestUDPPacket(t *testing.T) {
	h := ipv6.Header{
		Version:      6,
		TrafficClass: 0,
		FlowLabel:    0,
		PayloadLen:   0x001e,
		NextHeader:   17,
		HopLimit:     63,
		Src:          net.ParseIP("aaaa::c30c:0:0:7"),
		Dst:          net.ParseIP("aaaa::1"),
	}
	uHdr := UDPHeader{
		DstPort: 0xb25f,
		SrcPort: 0x1633,
		Length:  0x001e,
		Chksum:  0,
		Payload: []byte{0x60, 0x45, 0x22, 0xa6, 0x41, 0x0c, 0x80, 0xb1, 0x02, 0xff, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21},
	}

	cases := []struct {
		header ipv6.Header
		udpHdr UDPHeader
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
		//fmt.Printf("Succesful:: Case: %v, result: %x\n", c, b)

	}
}

//TestIpUDPPacket Test the full range of UDP to IP to bytes working
func TestIPUDPPacket(t *testing.T) {
	h := ipv6.Header{
		Version:      6,
		TrafficClass: 0,
		FlowLabel:    0,
		PayloadLen:   0x001e,
		NextHeader:   17,
		HopLimit:     63,
		Src:          net.ParseIP("aaaa::c30c:0:0:7"),
		Dst:          net.ParseIP("aaaa::1"),
	}
	uHdr := UDPHeader{
		DstPort: 0xb25f,
		SrcPort: 0x1633,
		Length:  0x001e,
		Chksum:  0,
		Payload: []byte{0x60, 0x45, 0x22, 0xa6, 0x41, 0x0c, 0x80, 0xb1, 0x02, 0xff, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21},
	}

	cases := []struct {
		header ipv6.Header
		udpHdr UDPHeader
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

		fullIP, err := Marschal(c.header, b)
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
		//fmt.Printf("Succesful:: Case: %v, result: %x\n", c, fullIP)

	}
}
