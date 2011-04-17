package main

const (
	sessionNs = "urn:ietf:params:xml:ns:xmpp-session"
	discoInfoNs = "http://jabber.org/protocol/disco#info"
	rosterNs = "jabber:iq:roster"
)

func (c *C2SConn) iq() {
	/* from, to := c.Attributes2("from", "to")
	if from != c.jid.String() {
		c.stanza.SetAttribute("from", c.jid.String())
	}

	if to == "" {
		c.SetAttr("to", c.jid.BareString())
	} */

	// c.route()

	id := c.MustAttr("id")

	c.MustToChild()

	xmlns := c.MustAttr("xmlns")
	switch xmlns {
	case sessionNs:
		c.Element("iq", "id", id, "type", "result").End()
	case discoInfoNs:
		c.StartElement("iq", "id", id, "type", "result").
			Element("query", "xmlns", discoInfoNs).
//			Element("identity", "category", "server", "type", "im", "name", "gxmppd")
			End()
	case rosterNs:
		c.StartElement("iq", "id", id, "type", "result").
			Element("query", "xmlns", rosterNs).
			End()
	default:
		c.StartElement("iq", "id", id, "type", "error").
			StartElement("error", "type", "cancel").
//				Element("feature-not-implemented", "xmlns", "urn:ietf:params:xml:ns:xmpp-stanzas").
				Element("service-unavailable", "xmlns", "urn:ietf:params:xml:ns:xmpp-stanzas").
			End()
	}
}
