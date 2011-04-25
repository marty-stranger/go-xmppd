package main

import (
//	"g/xml"
)

const stanzaErrorNs = "urn:ietf:params:xml:ns:xmpp-stanzas"

func (s *C2SStream) stanza() {
	defer func() {
		switch error := recover().(type) {
		case nil:
		/* case StanzaError:
			// hey maybe just use ret value ?
			stanza := error.Stanza
			stanza.Swap()
			stanza.Type = ErrorType
			stanza.Fragment = xml.NewBuilder().
				StartElement("error", "type", error.Type).
					Element(error.Condition, "xmlns", stanzaErrorNs).
					End()
			s.WriteStanza(stanza) */
		default:
			panic(error)
		}
	}()

	fragment := s.ReadElement()

	if fragment == nil { // should be </stream:stream>
		panic(nil)
		s.EndElementNoTrack("stream:stream")
		s.Close()
		return
	}

	stanza := newStanza(fragment)

	packet := &Packet{}
	packet.Src = s.To

	switch stanza.Kind {
	case IQKind:
		if stanza.To.Full == "" {
			packet.Dest = s.To.BareJid()
		} else {
			packet.Dest = stanza.To
		}
	case PresenceKind:
		// TODO check stanza.From

		// presence subscribe should be marked with bare, but broadcast with full, WTF?
		stanza.From = s.To.BareJid()
		packet.Dest = s.To.BareJid()
	case MessageKind:
		packet.Dest = stanza.To

		stanza.From = s.To
	}

	/* if stanza.To.Full != "" {
		packet.Dest = stanza.To
	} else {
		packet.Dest = s.jid.BareJid()
	} */
	packet.Stanza = stanza

	router.Ch <- packet
}
