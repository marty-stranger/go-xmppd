package main

type Router struct {
	ch	chan *Packet
}

var router = &Router{
	ch:	make(chan *Packet)}

func (r *Router) run() {
	for packet := range r.ch {
		debugln(packet)

		dest := packet.Dest
		if dest.Domain == serverName {
			if dest.Local != "" {
				if dest.Resource != "" {
					c2s.Ch <- packet
				} else {
					sm.Ch <- packet
				}
			} else {
				local.Ch <- packet
			}
		} else {
			s2s.Ch <- packet
		}

		// TODO components
	}
}

func init() {
	go router.run()
}

