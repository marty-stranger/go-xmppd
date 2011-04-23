package main

import (
	"os/signal"
)

func main() {
	debugln("started")
	for {
		switch (<-signal.Incoming).(signal.UnixSignal) {
		case signal.SIGINT:
			return
		}
	}
}
