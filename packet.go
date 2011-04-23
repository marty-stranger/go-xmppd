package main

import (
	"fmt"
)

type Packet struct {
	Src, Dest	Jid

	*Stanza
}

func (p *Packet) String() string {
	return fmt.Sprintf("Src:%s Dest:%s %s", p.Src.Full, p.Dest.Full, p.Stanza)
}

func (p *Packet) Swap() {
	p.Src, p.Dest = p.Dest, p.Src
	p.From, p.To = p.To, p.From
}
