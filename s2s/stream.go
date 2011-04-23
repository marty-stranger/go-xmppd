package main

import (
	"fmt"
	"net"
)

func address(domain string) string {
	_, addrs, err := net.LookupSRV("xmpp-server", "tcp", domain)

	if err != nil {
		return domain + ":5269"
	}

	// TODO choose one from addrs
	addr := addrs[0]
	return fmt.Sprint(addr.Target, ":", addr.Port)
}

type FromTo struct { From, To string }

type S2SStream struct {
	*Stream

	streamId	string
	pending		[]*Stanza
	// verified	map[FromTo]bool
	verified	bool
}

func newS2SStream() *S2SStream {
	return &S2SStream{}
}

func (s *S2SStream) connect(to string) {
	addr := address(to)

	conn, err := net.Dial("tcp", "", addr)
	if err != nil { panic(err) } // TODO handle error

	s.Stream = newStream(conn)
	s.To = makeJid(to)

	// NOTE xmlns:db is required for gmail.com, invalid-namespace otherwise
	version := "1.0"
	s.StartElement("stream:stream",
			"from", serverName,
			"to", s.To.Full,
			"version", version,
			"xmlns", "jabber:server",
			"xmlns:stream", streamNs,
			"xmlns:db", "jabber:server:dialback").
		Send()

	cursor := s.ReadStartElement().Cursor()

	s.streamId = cursor.MustAttr("id")

	version, _ = cursor.Attr("version")
	if version == "" {
		s.dialback()
	}

}

func (s *S2SStream) WriteStanza(stanza *Stanza) {
	if s.verified {
		s.Stream.WriteStanza(stanza)
	} else {
		s.pending = append(s.pending, stanza)
	}
}

func (s *S2SStream) sendPending() {
	for _, stanza := range s.pending {
		s.WriteStanza(stanza)
	}
	s.pending = nil
}
