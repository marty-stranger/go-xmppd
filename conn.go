package main

import (
	"net"

	"g/xml"
)

const (
	jabberClientNs = "jabber:client"
	streamNs = "http://etherx.jabber.org/streams"
)


type Conn struct {
	net.Conn
	*xml.Reader
	*xml.Writer
	*xml.Cursor
	// *Stanza

	jid	Jid

	connected	bool	// NOTE move them in flags?
	interested	bool
	available	bool
}

func newConn(netConn net.Conn) *Conn {
	return &Conn{
		Conn: netConn,
		Reader: xml.NewReader(&LogReadWriter{netConn}),
		Writer: xml.NewWriter(&LogReadWriter{netConn})}
}

func (c *Conn) stream() {
	cursor := c.ReadStartElement().Cursor()

	// TODO work with from attribute
	// streamTo, _ := cursor.Attr("from")
	version := cursor.MustAttr("version")
	if version != "1.0" { version = "" }

//	c.StartDocument().
	println("stream")
	c.StartElement("stream:stream", "from", serverName, "version", version,
			"xmlns", jabberClientNs, "xmlns:stream", streamNs).
		Send()
}

/* func (c *Conn) stanza() {
	// c.Stanza = Stanza{c.ReadElement().Cursor()}
} */

func (c *Conn) readElement() {
	c.Cursor = c.ReadElement().Cursor()
}

func (c *Conn) run() {
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
