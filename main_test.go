package sixlowpan

import (
	"context"
	"testing"
)

func TestMain(t *testing.T) {
	cases := []struct {
	}{
		{},
	}
	for _, c := range cases {
		config := Config{
			DebugLevel: DebugPacket,
			PortName:   "/dev/ttyUSB1",
		}
		ctx, cancel := context.WithCancel(context.Background())
		channel := make(chan []byte, 5)
		var got int
		go func() {
			got = Start(ctx, config, channel)
			if got != 0 {
				cancel()
			}
		}()
		select {
		case <-channel:
			cancel()
		case <-ctx.Done():

		}
		if got != 0 {
			t.Errorf("Did not exit properly %v, input: %v", got, c)
		}
	}
}
