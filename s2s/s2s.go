package main

type S2S struct {
	ch	chan *Packet
}

var s2s = S2S{
	ch:	make(chan *Packet)}

func (m *S2S) run() {
}
