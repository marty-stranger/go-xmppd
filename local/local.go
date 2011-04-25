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
		switch packet.Kind {
		case IQKind:
			m.iq(packet)
		}
	}
}

func (m *Local) iq(packet *Packet) {
	cursor := packet.Cursor()
	xmlns := cursor.MustAttr("xmlns")
	switch xmlns {
	case discoInfoNs:
		if packet.Type == GetType {
			fragment := xml.NewBuilder().
				Element("query", "xmlns", discoInfoNs).
				End()
			packet.Swap()
			packet.Type = ResultType
			packet.Fragment = fragment
			router.Ch <- packet
		}
	default:
		packet.Swap()
		packet.Type = ErrorType
		packet.Fragment = xml.NewBuilder().
			StartElement("error", "type", "cancel").
				Element("service-unavailable", "xmlns", "urn:ietf:params:xml:ns:xmpp-stanzas").
			End()
		router.Ch <- packet
	}
}

var local = &Local{
	Ch:	make(chan *Packet),
}

func init() {
	go local.run()
}
