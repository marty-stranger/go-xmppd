package main

const tlsRequired = true

func (s *C2SStream) header() {
	cursor := s.ReadStartElement().Cursor()

	from, _ := cursor.Attr("from")
	version := cursor.MustAttr("version")

	s.generateStreamId()

	s.StartDocument().
		StartElement("stream:stream",
			"from", serverName,
			"id", s.Id,
			"to", from,
			"version", "1.0",
			"xml:lang", "en",
			"xmlns", "jabber:client",
			"xmlns:stream", streamNs).
		Send()

	if version != "1.0" {
		s.streamError("unsupported-version")
	}
}

func (s *C2SStream) negotiate() {
	// TODO make configurable
	s.header()

	s.StartElement("stream:features").
		Element("starttls", "xmlns", tlsNs).
		End()

	cursor := s.ReadElement().Cursor()

	s.tls(cursor)

	s.header()
	s.StartElement("stream:features").
		StartElement("mechanisms", "xmlns", saslNs).
			Element("mechanism", "PLAIN").
		End()

	cursor = s.ReadElement().Cursor()

	local := s.sasl(cursor)
	s.header()

	s.StartElement("stream:features").
		Element("bind", "xmlns", bindNs).
//		Element("session", "xmlns", "urn:ietf:params:xml:ns:xmpp-session").
		End()

	cursor = s.ReadElement().Cursor()

	s.bind(local, cursor)
}
