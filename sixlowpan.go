//Package sixlowpan Package interface of sixlowpan for hecomm-fog
package sixlowpan

import (
	"fmt"
	"io"
	"log"

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
				fmt.Printf(string(packet))
			}

		//No special first character --> receiving payload packets
		default:
			if com.config.DebugLevel >= DebugPacket {
				fmt.Printf("SLIP Packet: %v, isPrefix: %v\n", packet, isPrefix)
			}
			if len(buf) < len(packet) {
				return 0, fmt.Errorf("Buf to small")
			}
			copy(buf, packet)
			return len(packet), nil
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

//Close Close the connection
func (com SLIP) Close() error {
	err := com.serial.Close()
	return err
}
