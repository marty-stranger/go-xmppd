package main

import (
	"g/xml"
)

func (m *SM) bareIQ(packet *Packet) {
	cursor := packet.Fragment.Cursor()

	switch cursor.MustAttr("xmlns") {
	// case discoInfoNs: m.discoInfoIQ()
	case rosterNs:	m.rosterIQ(packet)
	default:
		packet.Swap()
		packet.Type = "error"
		packet.Fragment = xml.NewBuilder().
			StartElement("error", "type", "cancel").
				Element("service-unavailable", "xmlns", "urn:ietf:params:xml:ns:xmpp-stanzas").
			End()
		router.ch <- packet
	}
}

func (m *SM) iq(packet *Packet) {
	if packet.Type == "result" {
	} else {
		if packet.To.Resource == "" {
			m.bareIQ(packet)
		} else {
			// m.iqFullTo()
		}
	}
}

/* func (m *SMOut) resultIQ(fragment *xml.Fragment) {
	stanza := &Stanza{
		Name: "iq",
		From: m.To,
		Id: m.Id,
		To: m.From,
		Type: "result"}
	stanza.Fragment = fragmnent
	m.route(stanza)
} */

