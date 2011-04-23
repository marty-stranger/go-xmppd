package main

import (
	"net"
)

type S2S struct {
	Ch	chan *Packet
	streams	map[string]*S2SStream
}

var s2s = &S2S{
	Ch:		make(chan *Packet),
	streams:	make(map[string]*S2SStream)}

func (m *S2S) run() {
	for packet := range m.Ch {
		debugln(packet)
		domain := packet.To.Domain

		stream := m.streams[domain]
		if stream == nil {
			stream = newS2SStream()
			go stream.connect(domain)
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
