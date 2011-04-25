package main

type Router struct {
	Ch	chan *Packet
}

// NOTE does 1000 chan length eliminate deadlock problem ?
var router = &Router{
	Ch:	make(chan *Packet, 1000),
}

func (r *Router) run() {
	for packet := range r.Ch {
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

