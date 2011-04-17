package s2s

import (
	"crypto/rand"
	"fmt"
	"net"

	"g/xmppd/base"
)

func Run() {
	l, e := net.Listen("tcp", "0.0.0.0:5269")
	if e != nil { panic(e) }

	for {
		c, e := l.Accept()
		if e != nil { panic(e) }

		s := &Stream{Stream: base.NewStream(c)}
		go s.accept()
	}
}

func address(domain string) string {
	_, addrs, err := net.LookupSRV("xmpp-server", "tcp", domain)

	if err != nil {
		return domain + ":5269"
	}

	// TODO choose one from addrs
	addr := addrs[0]
	return fmt.Sprint(addr.Target, ":", addr.Port)
}

type Stream struct {
	*base.Stream

	streamTo	string
	streamId	string
}

func Connect(to string) {
	addr := address(to)

	conn, err := net.Dial("tcp", "", addr)
	if err != nil { panic(err) }

	stream := &Stream{
		Stream: base.NewStream(conn),
		streamTo: to}
	stream.connect()
}

func (s *Stream) connect() {
	// NOTE xmlns:db is required for gmail.com, invalid-namespace otherwise
	version := "1.0"
	s.StartElement("stream:stream",
			"from", base.ServerName,
			"to", s.streamTo,
			"version", version,
			"xmlns", "jabber:server",
			"xmlns:stream", base.StreamNs,
			"xmlns:db", "jabber:server:dialback").
		Send()

	cursor := s.ReadStartElement().Cursor()

	s.streamId = cursor.MustAttr("id")

	version, _ = cursor.Attr("version")
	if version == "" {
		s.dialback()
	}

}

func (s *Stream) accept() {
	cursor := s.ReadStartElement().Cursor()

	from, _ := cursor.Attr("from")
	version, _ := cursor.Attr("version")
	dialback, _ := cursor.Attr("xmlns:db")

	bytes := make([]byte, 16)
	rand.Read(bytes)

	id := fmt.Sprintf("%x", bytes)

	s.StartElement("stream:stream",
			"from", base.ServerName,
			"id", id,
			"to", from,
			"version", version,
			"xmlns", "jabber:server",
			"xmlns:stream", base.StreamNs,
			"xmlns:db", dialback).
		Send()

	s.dialbackAccept()
}
