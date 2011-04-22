package main

type C2SStream struct {
	*Stream

	jid	Jid

	connected	bool	// NOTE move them in flags?
	interested	bool
	available	bool
}

func (s *C2SStream) stream() {
	cursor := s.ReadStartElement().Cursor()

	// TODO work with from attribute
	// streamTo, _ := cursor.Attr("from")
	version := cursor.MustAttr("version")
	if version != "1.0" { version = "" }

	s.StartDocument().
		StartElement("stream:stream", "from", serverName, "version", version,
			"xmlns", "jabber:client", "xmlns:stream", streamNs).
		Send()
}

func (s *C2SStream) run() {
	// TODO unify init logic
	s.stream()

	s.StartElement("stream:features").
		Element("starttls", "xmlns", "urn:ietf:params:xml:ns:xmpp-tls").
		End()

	cursor := s.ReadElement().Cursor()

	s.tls(cursor)

	s.stream()
	s.StartElement("stream:features").
		StartElement("mechanisms", "xmlns", saslNs).
			Element("mechanism", "PLAIN").
		End()

	cursor = s.ReadElement().Cursor()

	local := s.sasl(cursor)
	s.stream()

	s.StartElement("stream:features").
		Element("bind", "xmlns", bindNs).
//		Element("session", "xmlns", "urn:ietf:params:xml:ns:xmpp-session").
		End()

	cursor = s.ReadElement().Cursor()
	s.bind(local, cursor)

	for {
		stanza := newStanza(s.ReadElement())

		packet := &Packet{}
		packet.Src = s.jid

		switch stanza.Name {
		case "iq":
			if stanza.To.Full == "" {
				packet.Dest = s.jid.BareJid()
			} else {
				packet.Dest = stanza.To
			}
		case "presence":
			packet.Dest = s.jid.BareJid()
		case "message":
			packet.Dest = stanza.To

			stanza.From = s.jid
		}

		/* if stanza.To.Full != "" {
			packet.Dest = stanza.To
		} else {
			packet.Dest = s.jid.BareJid()
		} */
		packet.Stanza = stanza

		router.ch <- packet
	}
}
