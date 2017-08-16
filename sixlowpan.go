//Package sixlowpan Package interface of sixlowpan for hecomm-fog
package sixlowpan

import (
	"fmt"
	"io"
	"log"

	"golang.org/x/net/ipv6"

	"github.com/Lobaro/slip"
	"github.com/jacobsa/go-serial/serial"
)

const (
	//DebugNone No printing
	DebugNone uint8 = iota
	//DebugPacket Print packet
	DebugPacket
	//DebugAll Print everything!
	DebugAll
)

//Config Configuration of interface
type Config struct {
	DebugLevel uint8
	PortName   string
}

//SLIP The struct of sixlowpan slip node, io readwritecloser
type SLIP struct {
	config    Config
	serial    io.ReadWriteCloser
	slipread  *slip.Reader
	slipwrite *slip.Writer
}

//Open interface up to sixlowpan SLIP
func Open(config Config) (com io.ReadWriteCloser, err error) {
	log.Println("Opening serial hecomm SLIP port...")
	sl := SLIP{}
	sl.config = config
	//Serial setup
	//To use github.com/jacobsa/go-serial/serial or go.bug.st/serial.v1
	options := serial.OpenOptions{
		PortName:        config.PortName,
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}

	port, err := serial.Open(options)
	if err != nil {
		return com, err
	}
	sl.serial = port

	//reader := slip.NewReader(port)
	sl.slipread = slip.NewReader(port)
	sl.slipwrite = slip.NewWriter(port)

	return sl, err
}

//UdpHeaderLen The length of a standard UDP header is 8 bytes
const UdpHeaderLen = 8

//Read Read until next packet received
func (com SLIP) Read(buf []byte) (n int, err error) {
	for {
		packet, isPrefix, err := com.slipread.ReadPacket()
		if err != nil {
			return 0, err
		}
		switch packet[0] {
		//Case of \r --> receiving debug lines through slip
		case 0x0D:
			if com.config.DebugLevel >= DebugAll {
				fmt.Printf("Debug packet: %v", string(packet))
			}

		//No special first character --> receiving payload packets
		default:
			if com.config.DebugLevel >= DebugPacket {
				if len(packet) < (ipv6.HeaderLen + UdpHeaderLen) {
					log.Fatalln("Packet to small to fit ipv6 + udp header!!")
				}
				fmt.Printf("SLIP Packet payload: %v, isPrefix: %v\n", packet[ipv6.HeaderLen+UdpHeaderLen:], isPrefix)
				header, err := ipv6.ParseHeader(packet[:ipv6.HeaderLen])
				if err != nil {
					fmt.Printf("Error in processing ipv6 header\n")
				} else {
					fmt.Printf("Header %v\n", header.String())
				}
			}
			if len(buf) < len(packet) {
				return 0, fmt.Errorf("Buf to small")
			}
			copy(buf, packet[ipv6.HeaderLen+UdpHeaderLen:])
			return len(packet) - (ipv6.HeaderLen + UdpHeaderLen), nil
		}
	}
}

//Write	Write a packet over sixlowpan SLIP
func (com SLIP) Write(p []byte) (n int, err error) {
	err = com.slipwrite.WritePacket(p)
	if err != nil {
		return 0, err
	}
	return len(p), err
}

//Marschal Compile the ipv6 packet
func Marschal(h ipv6.Header, p []byte) (b []byte, err error) {
	if h.PayloadLen != len(p) {
		return b, fmt.Errorf("PayloadLen and payload parameter does not match")
	}
	b = make([]byte, (ipv6.HeaderLen + h.PayloadLen))
	//Big-endian or middle endianess
	b[0] = byte(h.Version)<<4 | byte(h.TrafficClass)>>4 //version to lower, trafficclass only highbits
	b[1] = byte(h.TrafficClass)<<4 | byte((h.FlowLabel >> 16))
	b[2] = byte(h.FlowLabel >> 8)
	b[3] = byte(h.FlowLabel)
	b[4] = byte(h.PayloadLen >> 8)
	b[5] = byte(h.PayloadLen)
	b[6] = byte(h.NextHeader)
	b[7] = byte(h.HopLimit)
	copy(b[8:24], h.Src)
	copy(b[24:40], h.Dst)

	copy(b[40:40+h.PayloadLen], p)
	return b, nil
	/* h := &Header{
		Version:      int(b[0]) >> 4,
		TrafficClass: int(b[0]&0x0f)<<4 | int(b[1])>>4,
		FlowLabel:    int(b[1]&0x0f)<<16 | int(b[2])<<8 | int(b[3]),
		PayloadLen:   int(binary.BigEndian.Uint16(b[4:6])),
		NextHeader:   int(b[6]),
		HopLimit:     int(b[7]),
	}
	h.Src = make(net.IP, net.IPv6len)
	copy(h.Src, b[8:24])
	h.Dst = make(net.IP, net.IPv6len)
	copy(h.Dst, b[24:40]) */
}

//Close Close the connection
func (com SLIP) Close() error {
	err := com.serial.Close()
	return err
}
