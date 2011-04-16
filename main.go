package main

import (
	"net"
)

const (
	serverName = "gxmppd.org"
)

func main() {
	l, e := net.Listen("tcp", "0.0.0.0:5222")
	if e != nil { panic(e) }

	for {
		c, e := l.Accept()
		// TODO just report about error
		if e != nil { panic(e) }

		go newConn(c).run()
	}
}

