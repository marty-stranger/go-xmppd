package main

import (
	"sync"
)

type Conns struct {
	conns map[string]map[string]*Conn
	mutex sync.RWMutex
}

func (c *Conns) Add(local, resource string, conn *Conn) {
	localConns := c.conns[local]
	if localConns == nil {
		localConns = make(map[string]*Conn)
		c.conns[local] = localConns
	}

	localConns[resource] = conn
}

func (c *Conns) Available(local string) []*Conn {
	conns := []*Conn{}
	for _, conn := range c.conns[local] {
		if conn.available {
			conns = append(conns, conn)
		}
	}

	return conns
}

var conns = &Conns{conns: make(map[string]map[string]*Conn)}


