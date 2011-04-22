package main

type C2SConn struct {
	*Conn

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

	c.StartDocument().
		StartElement("stream:stream", "from", serverName, "version", version,
			"xmlns", "jabber:client", "xmlns:stream", streamNs).
		Send()
}

func (c *C2SConn) run() {
	// TODO unify init logic
	c.stream()

	c.StartElement("stream:features").
		Element("starttls", "xmlns", "urn:ietf:params:xml:ns:xmpp-tls").
		End()

	cursor := c.ReadElement().Cursor()

	c.tls(cursor)

	c.stream()
	c.StartElement("stream:features").
		StartElement("mechanisms", "xmlns", saslNs).
			Element("mechanism", "PLAIN").
		End()

	cursor = c.ReadElement().Cursor()

	local := c.sasl(cursor)
	c.stream()

	c.StartElement("stream:features").
		Element("bind", "xmlns", bindNs).
//		Element("session", "xmlns", "urn:ietf:params:xml:ns:xmpp-session").
		End()

	cursor = c.ReadElement().Cursor()
	c.bind(local, cursor)

	for {
		stanza := newStanza(c.ReadElement())

		packet := &Packet{}
		packet.Src = c.jid

		if stanza.To.Full != "" {
			packet.Dest = stanza.To
		} else {
			packet.Dest = c.jid.BareJid()
		}
		packet.Stanza = stanza

		router.ch <- packet
	}
}

func (c *C2SConn) writeStanza(stanza *Stanza) {
	c.StartElement(stanza.Name, "from", stanza.From.Full, "id", stanza.Id,
		"to", stanza.To.Full, "type", stanza.Type).
		Raw(stanza.Fragment.String()).
		End()
}
