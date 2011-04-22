package main

import (
	"fmt"
	"net"
)

func runS2S() {
	l, e := net.Listen("tcp", "0.0.0.0:5269")
	if e != nil { panic(e) }

	for {
		c, e := l.Accept()
		if e != nil { panic(e) }

		s := &S2SConn{Conn: newConn(c)}
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

type FromTo struct { From, To string }

type S2SConn struct {
	*Conn

	streamTo	string
	streamId	string
	pending		[]*Stanza
	// verified	map[FromTo]bool
	verified	bool
}

func newS2SConn(to string) *S2SConn {
	return &S2SConn{
		streamTo: to}
//		verified: make(map[FromTo]bool)}
}

func (c *S2SConn) connect() {
	addr := address(c.streamTo)

	conn, err := net.Dial("tcp", "", addr)
	if err != nil { panic(err) } // TODO handle error

	c.Conn = newConn(conn)

	// NOTE xmlns:db is required for gmail.com, invalid-namespace otherwise
	version := "1.0"
	c.StartElement("stream:stream",
			"from", serverName,
			"to", c.streamTo,
			"version", version,
			"xmlns", "jabber:server",
			"xmlns:stream", streamNs,
			"xmlns:db", "jabber:server:dialback").
		Send()

	cursor := c.ReadStartElement().Cursor()

	c.streamId = cursor.MustAttr("id")

	version, _ = cursor.Attr("version")
	if version == "" {
		c.dialback()
	}

}

func (c *S2SConn) sendStanza(stanza *Stanza) {
	if c.verified {
		c.StartElement(stanza.Name, "from", stanza.From.Full, "id", stanza.Id,
				"to", stanza.To.Full, "type", stanza.Type).
			Raw(stanza.Fragment.String()).
			End()
	} else {
		c.pending = append(c.pending, stanza)
	}
}

func (c *S2SConn) sendPending() {
	for _, stanza := range c.pending {
		c.sendStanza(stanza)
	}
	c.pending = nil
}
