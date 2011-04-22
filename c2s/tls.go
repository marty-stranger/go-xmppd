package main

import (
	"crypto/tls"

	"g/xml"
)

const tlsNs = "urn:ietf:params:xml:ns:xmpp-tls"

func (s *C2SStream) tls(cursor *xml.Cursor) {
	s.Element("proceed", "xmlns", tlsNs).End()

	// TODO add certificates
	cert, e := tls.LoadX509KeyPair("etc/tls/server.crt", "etc/tls/server.key")
	if e != nil { panic(e) }

	config := &tls.Config{}
	config.Certificates = []tls.Certificate{cert}

	tlsConn := tls.Server(s.Conn, config)
	e = tlsConn.Handshake()
	if e != nil { panic(e) }

	s.Stream = newStream(tlsConn)
}
