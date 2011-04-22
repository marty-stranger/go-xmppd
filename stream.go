package main

import (
	"net"

	"g/xml"
)

const (
	streamNs = "http://etherx.jabber.org/streams"
)

type Stream struct {
	net.Conn
	*xml.Reader
	*xml.Writer
}

func newStream(netConn net.Conn) *Stream {
	return &Stream{
		netConn,
		xml.NewReader(&LogReadWriter{netConn}),
		xml.NewWriter(&LogReadWriter{netConn})}
}


func (c *Stream) WriteStanza(stanza *Stanza) {
	c.StartElement(stanza.Name, "from", stanza.From.Full, "id", stanza.Id,
		"to", stanza.To.Full, "type", stanza.Type).
		Raw(stanza.Fragment.String()).
		End()
}
