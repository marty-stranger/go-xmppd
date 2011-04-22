package main

var s2sConns = make(map[string]*S2SConn)

func routeS2S(stanza *Stanza) {
	domain := stanza.To.Domain

	conn := s2sConns[domain]
	if conn == nil {
		conn := newS2SConn(domain)
		go conn.connect()
		s2sConns[domain] = conn
	}

	conn.sendStanza(stanza)
}
