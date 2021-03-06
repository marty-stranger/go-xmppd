package main

import (
	"github.com/pavelrosputko/go-xml"
)

func (p SMPacket) bareIQ() {
	debugln("")
	cursor := p.Cursor()

	switch cursor.MustAttr("xmlns") {
	// case discoInfoNs: m.discoInfoIQ()
	case rosterNs: p.rosterIQ()
	default:
		p.Swap()
		p.Type = ErrorType
		p.Fragment = xml.NewBuilder().
			StartElement("error", "type", "cancel").
				Element("service-unavailable", "xmlns", "urn:ietf:params:xml:ns:xmpp-stanzas").
			End()
		router.Ch <- p.Packet
	}
}

func (p SMPacket) iq() {
	if p.Type == ResultType {
	} else {
		if p.To.Resource == "" {
			p.bareIQ()
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

