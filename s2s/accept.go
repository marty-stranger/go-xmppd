package main

import (
	"crypto/rand"
	"fmt"
)

func (s *S2SStream) accept() {
	cursor := s.ReadStartElement().Cursor()

	from, _ := cursor.Attr("from")
	version, _ := cursor.Attr("version")
	dialback, _ := cursor.Attr("xmlns:db")

	bytes := make([]byte, 16)
	rand.Read(bytes)

	id := fmt.Sprintf("%x", bytes)

	s.StartElement("stream:stream",
			"from", serverName,
			"id", id,
			"to", from,
			"version", version,
			"xmlns", "jabber:server",
			"xmlns:stream", streamNs,
			"xmlns:db", dialback).
		Send()

	s.dialbackAccept()

	for {
		cursor := s.ReadElement().Cursor()

		if cursor.Name() == "db:result" {
			continue
		}

		stanza := newStanza(cursor.Fragment)

		packet := &Packet{}
		packet.Dest = stanza.To
		packet.Stanza = stanza
		router.ch <- packet
	}
}

