package main

import (
	"fmt"
)

func init() {
	go sm.run()
}

type SM struct {
	ch		chan *Packet
	id		int
	connected	map[string]map[string]bool
	available	map[string]map[string]bool
}

var sm = SM{
	ch:		make(chan *Packet),
	connected:	make(map[string]map[string]bool),
	available:	make(map[string]map[string]bool)}

func (m *SM) run() {
	for packet := range m.ch {
		switch packet.Name {
		case "iq": m.iq(packet)
		case "presence": m.presence(packet)
		// case "message": m.message(packet)
		}
	}
}

func (m *SM) nextId() string {
	m.id++
	return fmt.Sprint(m.id)
}

func (m *SM) BindResource(local, resource string) bool {
	resources := m.connected[local]
	if resources == nil {
		resources = make(map[string]bool)
		m.connected[local] = resources
	}

	if resources[resource] {
		return false
	}

	resources[resource] = true
	return true

}

