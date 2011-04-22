package main

import (
	"os/signal"
)

func main() {
	go runC2S()
	go runS2S()

	for {
		switch (<-signal.Incoming).(signal.UnixSignal) {
		case signal.SIGINT:
			return
		}
	}
}
