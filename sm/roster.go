package main

import (
	"g/xml"
)

const rosterNs = "jabber:iq:roster"

func (m *SM) rosterIQ(packet *Packet) {
	switch packet.Type {
	case "get":
		fragment := xml.NewBuilder().
			Element("query", "xmlns", rosterNs).
			End()
		packet.Swap()
		packet.Type = "result"
		packet.Fragment = fragment
		router.ch <- packet
	case "set":
	}
}
