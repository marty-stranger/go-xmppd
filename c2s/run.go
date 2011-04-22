package main

import (
	"net"
)

func runC2S() {
	l, e := net.Listen("tcp", "0.0.0.0:5222")
	if e != nil { panic(e) }

	for {
		c, e := l.Accept()
		// TODO just report about error
		if e != nil { panic(e) }

		c2sConn := &C2SConn{Conn: newConn(c)}
		go c2sConn.run()
	}
}
