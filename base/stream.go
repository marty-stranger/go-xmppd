package base

import (
	"net"

	"g/xml"
)

const (
	StreamNs = "http://etherx.jabber.org/streams"
)

type Stream struct {
	net.Conn
	*xml.Reader
	*xml.Writer
}

func NewStream(netConn net.Conn) *Stream {
	return &Stream{
		netConn,
		xml.NewReader(&LogReadWriter{netConn}),
		xml.NewWriter(&LogReadWriter{netConn})}
}

