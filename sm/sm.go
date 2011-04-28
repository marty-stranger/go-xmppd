package main

import (
	"fmt"

	"github.com/pavelrosputko/go-xml"
)

func init() {
	go sm.run()
}

type SM struct {
	Ch		chan *Packet
	Sessions	map[string]map[string]*Session

	id		int
}

type Session struct {
	Jid		Jid
	Interested	bool
	Available	bool
	Presence	*xml.Fragment
}

type SMPacket struct {
	*Packet
}

func (sm *SM) GetSession(user, resource string) *Session {
	if sessions := sm.Sessions[user]; sessions != nil {
		return sessions[resource]
	}

	return nil
}

func (sm *SM) HasAvailable(user string) bool {
	for _, session := range sm.Sessions[user] {
		if session.Available { return true }
	}

	return false
}

var sm = SM{
	Ch:		make(chan *Packet),
	Sessions:	make(map[string]map[string]*Session)}

func (m *SM) run() {
	for packet := range m.Ch {
		debugln(packet)
		smPacket := SMPacket{packet}
		switch packet.Kind {
		case IQKind: smPacket.iq()
		case PresenceKind: smPacket.presence()
		// case "message": m.message(packet)
		}
	}
}

func (m *SM) nextId() string {
	m.id++
	return fmt.Sprint(m.id)
}

func (m *SM) BindResource(user, resource string) bool {
	sessions := m.Sessions[user]
	if sessions == nil {
		sessions = make(map[string]*Session)
		m.Sessions[user] = sessions
	}

	if sessions[resource] != nil {
		return false
	}

	sessions[resource] = &Session{
		Jid: makeJid(user + "@" + serverName + "/" + resource),
	}

	return true
}

func (m *SM) UnbindResource(user, resource string) {
	// TODO sync
	m.Sessions[user][resource] = nil, false
}

func (s *Session) Unavailable() {
	s.Available = false
}
