package main

import (
	"crypto/tls"

	"g/xml"
)

const tlsNs = "urn:ietf:params:xml:ns:xmpp-tls"

func (c *C2SConn) tls(cursor *xml.Cursor) {
	c.Element("proceed", "xmlns", tlsNs).End()

	// TODO add certificates
	cert, e := tls.LoadX509KeyPair("etc/tls/server.crt", "etc/tls/server.key")
	if e != nil { panic(e) }

	config := &tls.Config{}
	config.Certificates = []tls.Certificate{cert}

	tlsConn := tls.Server(c.Conn, config)
	e = tlsConn.Handshake()
	if e != nil { panic(e) }

	c.Conn = newConn(tlsConn)
}
