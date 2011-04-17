package main

import (
	"sync"
)

type C2SConns struct {
	conns map[string]map[string]*C2SConn
	mutex sync.RWMutex
}

func (c *C2SConns) Add(local, resource string, conn *C2SConn) {
	localConns := c.conns[local]
	if localConns == nil {
		localConns = make(map[string]*C2SConn)
		c.conns[local] = localConns
	}

	localConns[resource] = conn
}

func (c *C2SConns) Available(local string) []*C2SConn {
	conns := []*C2SConn{}
	for _, conn := range c.conns[local] {
		if conn.available {
			conns = append(conns, conn)
		}
	}

	return conns
}

var c2sConns = &C2SConns{conns: make(map[string]map[string]*C2SConn)}
