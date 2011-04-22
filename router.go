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
		fmt.Println("router: packet =", packet)

		println("d =", packet.Dest.Full)
		if packet.Dest.Full == serverName {
			local.ch <- packet
			continue
		}

		// from c2s
		if packet.Src.Domain == serverName && packet.Src.Resource != "" {
			sm.ch <- packet
			continue
		}

		if packet.Dest.Domain == serverName {
			if packet.Dest.Resource == "" {
				sm.ch <- packet
			} else {
				c2s.ch <- packet
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

