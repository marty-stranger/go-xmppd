package main

import (
	"net"
	"os"

	"g/xml"
)

const (
	streamNs = "http://etherx.jabber.org/streams"
)

type Stream struct {
	net.Conn
	*xml.Reader
	*xml.Writer

	To	Jid
}

func newStream(netConn net.Conn) *Stream {
	s := &Stream{Conn: netConn}
	s.Reader = xml.NewReader(s)
	s.Writer = xml.NewWriter(s)
	return s
}


func (s *Stream) WriteStanza(stanza *Stanza) {
	s.StartElement(stanza.Name, "from", stanza.From.Full, "id", stanza.Id,
		"to", stanza.To.Full, "type", stanza.Type).
		Raw(stanza.Fragment.String()).
		End()
}

func (s *Stream) Read(b []byte) (int, os.Error) {
	n, e := s.Conn.Read(b)
	debugln(s.To, string(b))
	return n, e
}

func (s *Stream) Write(b []byte) (int, os.Error) {
	n, e := s.Conn.Write(b)
	debugln(s.To, string(b))
	return n, e
}
