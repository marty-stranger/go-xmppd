package main

import (
	"g/xml"
)

const discoInfoNs = "http://jabber.org/protocol/disco#info"

type Local struct {
	Ch	chan *Packet
}

func (m *Local) run() {
	for packet := range m.Ch {
		debugln(packet)
		switch packet.Name {
		case "iq":
			m.iq(packet)
		}
	}
}

func (m *Local) iq(packet *Packet) {
	if packet.Type == "get" {
		fragment := xml.NewBuilder().
			Element("query", "xmlns", discoInfoNs).
			End()
		packet.Swap()
		packet.Type = "result"
		packet.Fragment = fragment
		router.ch <- packet
	}
}

var local = &Local{
	Ch:	make(chan *Packet)}

func init() {
	go local.run()
}
