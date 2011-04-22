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
	connected	UsersResources
	available	UsersResources
	interested	UsersResources
}

type UsersResources map[string]map[string]bool

func (ur UsersResources) Add(user, resource string) {
	resources := ur[user]
	if resources == nil {
		resources = make(map[string]bool)
		ur[user] = resources
	}

	resources[resource] = true
}

var sm = SM{
	ch:		make(chan *Packet),
	connected:	make(UsersResources),
	available:	make(UsersResources),
	interested:	make(UsersResources)}

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

