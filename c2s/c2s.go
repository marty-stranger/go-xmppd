package main

import (
	"net"
	"fmt"
	"sync"
)

type C2S struct {
	ch	chan *Packet
	streams	map[string]*C2SStream
	sync.RWMutex
}

var c2s = &C2S{
	ch:		make(chan *Packet),
	streams:	make(map[string]*C2SStream)}

func (m *C2S) run() {
	for packet := range m.ch {
		fmt.Println("C2S#run", packet)

		stream := m.streams[packet.Dest.Local + "/" + packet.Dest.Resource]
		if stream != nil {
			stream.WriteStanza(packet.Stanza)
		} else {
			println("stream nil")
			// TODO route error back
		}
	}
}

func (m *C2S) Add(local, resource string, stream *C2SStream) {
	// TODO / is inefficient
	m.streams[local + "/" + resource] = stream
}

func runC2S() {
	l, e := net.Listen("tcp", "0.0.0.0:5222")
	if e != nil { panic(e) }

	for {
		c, e := l.Accept()
		// TODO just report about error
		if e != nil { panic(e) }

		s := &C2SStream{Stream: newStream(c)}
		go s.run()
	}
}

func init() {
	go runC2S()
	go c2s.run()
}

