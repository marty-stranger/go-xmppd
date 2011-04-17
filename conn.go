package main

import (
	"net"

	"g/xml"
)

const (
	streamNs = "http://etherx.jabber.org/streams"
)

type Conn struct {
	net.Conn
	*xml.Reader
	*xml.Writer
}

func newConn(netConn net.Conn) *Conn {
	return &Conn{
		netConn,
		xml.NewReader(&LogReadWriter{netConn}),
		xml.NewWriter(&LogReadWriter{netConn})}
}

