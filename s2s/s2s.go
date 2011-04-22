package main

import (
	"fmt"
	"net"
)

type S2S struct {
	ch	chan *Packet
	streams	map[string]*S2SStream
}

var s2s = &S2S{
	ch	:	make(chan *Packet),
	streams:	make(map[string]*S2SStream)}

func (m *S2S) run() {
	for packet := range m.ch {
		fmt.Println("s2s", packet)
		domain := packet.To.Domain

		stream := m.streams[domain]
		if stream == nil {
			stream = newS2SStream(domain)
			go stream.connect()
			m.streams[domain] = stream
		}

		stream.WriteStanza(packet.Stanza)
	}
}

func runS2S() {
	l, e := net.Listen("tcp", "0.0.0.0:5269")
	if e != nil { panic(e) }

	for {
		c, e := l.Accept()
		if e != nil { panic(e) }

		s := &S2SStream{Stream: newStream(c)}
		go s.accept()
	}
}

func init() {
	go s2s.run()
	go runS2S()
}
