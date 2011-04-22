package main

import (
	"fmt"
	"sync"
)

type C2S struct {
	ch	chan *Packet
	conns	map[string]*C2SConn
	sync.RWMutex
}

var c2s = &C2S{
	ch:	make(chan *Packet),
	conns:	make(map[string]*C2SConn)}

func init() {
	go c2s.run()
}

func (m *C2S) run() {
	for packet := range m.ch {
		fmt.Println("C2S#run", packet)

		conn := m.conns[packet.Dest.Local + "/" + packet.Dest.Resource]
		if conn != nil {
			conn.writeStanza(packet.Stanza)
		} else {
			println("conn nil")
			// TODO route error back
		}
	}
}

func (m *C2S) Add(local, resource string, conn *C2SConn) {
	// TODO / is inefficient
	m.conns[local + "/" + resource] = conn
}
