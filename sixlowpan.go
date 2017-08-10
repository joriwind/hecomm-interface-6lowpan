//Package sixlowpan Package interface of sixlowpan for hecomm-fog
package sixlowpan

import (
	"context"
	"fmt"
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

//Start start
func Start(ctx context.Context, config Config, packetChannel chan []byte) int {
	log.Println("Start of 6LoWPAN interface!")
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
		log.Fatalf("serial.Open: %v", err)
	}
	defer port.Close()

	reader := slip.NewReader(port)

	for {
		packet, isPrefix, err := reader.ReadPacket()
		if err != nil {
			log.Fatalf("6LoWPAN interface: Error in reading SLIP: %v", err)
		}
		switch packet[0] {
		//Case of \r --> receiving debug lines through slip
		case 0x0D:
			if config.DebugLevel >= DebugAll {
				fmt.Printf(string(packet))
			}

		//No special first character --> receiving payload packets
		default:
			if config.DebugLevel >= DebugPacket {
				fmt.Printf("SLIP Packet: %v, isPrefix: %v\n", packet, isPrefix)
			}
			packetChannel <- packet
		}

		select {
		case <-ctx.Done():
			return 0
		default:
		}
	}

}

func main() {
	fmt.Println("Hello")
}
