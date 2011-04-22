package main

import (
	"os/signal"
)

func main() {
	for {
		switch (<-signal.Incoming).(signal.UnixSignal) {
		case signal.SIGINT:
			return
		}
	}
}
