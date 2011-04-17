package main

import (
	"os/signal"
)

func main() {
	go RunC2S()
	go RunS2S()

	for {
		switch (<-signal.Incoming).(signal.UnixSignal) {
		case signal.SIGINT:
			return
		}
	}
}
