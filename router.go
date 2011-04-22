package main

import (
	"fmt"
)

type Router struct {
	ch	chan *Packet
}

var router = &Router{
	ch:	make(chan *Packet)}

func (r *Router) run() {
	for packet := range r.ch {
		fmt.Println("router", packet)

		dest := packet.Dest
		if dest.Domain == serverName {
			if dest.Local != "" {
				if dest.Resource != "" {
					c2s.ch <- packet
				} else {
					sm.ch <- packet
				}
			} else {
				local.ch <- packet
			}
		} else {
			s2s.ch <- packet
		}

		// TODO components
	}
}

func init() {
	go router.run()
}

