package main

import (
	"g/xml"
)

type C2SConn struct {
	*Conn

	*xml.Cursor

	jid	Jid

	connected	bool	// NOTE move them in flags?
	interested	bool
	available	bool
}

func (c *C2SConn) stream() {
	cursor := c.ReadStartElement().Cursor()

	// TODO work with from attribute
	// streamTo, _ := cursor.Attr("from")
	version := cursor.MustAttr("version")
	if version != "1.0" { version = "" }

//	c.StartDocument().
	println("stream")
	c.StartElement("stream:stream", "from", serverName, "version", version,
			"xmlns", "jabber:client", "xmlns:stream", streamNs).
		Send()
}

/* func (c *Conn) stanza() {
	// c.Stanza = Stanza{c.ReadElement().Cursor()}
} */

func (c *C2SConn) readElement() {
	c.Cursor = c.ReadElement().Cursor()
}

func (c *C2SConn) run() {
	// TODO unify init logic
	c.stream()

	c.StartElement("stream:features").
		Element("starttls", "xmlns", "urn:ietf:params:xml:ns:xmpp-tls").
		End()

	c.readElement()

	c.tls()

	c.stream()
	c.StartElement("stream:features").
		StartElement("mechanisms", "xmlns", saslNs).
			Element("mechanism", "PLAIN").
		End()

	c.readElement()

	local := c.sasl()
	c.stream()

	c.StartElement("stream:features").
		Element("bind", "xmlns", bindNs).
//		Element("session", "xmlns", "urn:ietf:params:xml:ns:xmpp-session").
		End()

	c.readElement()

	c.bind(local)

	for {
		// c.stanza()
		c.readElement()
		switch c.Name() {
		case "iq": c.iq()
		case "presence": c.presence()
		// case "message": c.message()
		}
	}
}
