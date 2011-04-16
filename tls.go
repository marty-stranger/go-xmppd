package main

import (
	"crypto/tls"

	"g/xml"
)

const tlsNs = "urn:ietf:params:xml:ns:xmpp-tls"

func (c *Conn) tls() {
	c.Element("proceed", "xmlns", tlsNs).End()

	// TODO add certificates
	cert, e := tls.LoadX509KeyPair("etc/tls/server.crt", "etc/tls/server.key")
	if e != nil { panic(e) }

	config := &tls.Config{}
	config.Certificates = []tls.Certificate{cert}

	tlsConn := tls.Server(c.Conn, config)
	e = tlsConn.Handshake()
	if e != nil { panic(e) }

	c.Conn = tlsConn
	c.Reader = xml.NewReader(&LogReadWriter{tlsConn})
	c.Writer = xml.NewWriter(&LogReadWriter{tlsConn})
}
