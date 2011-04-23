package main

import (
	"g/xml"
)

func (p SMPacket) presence() {
	if p.Src.Resource != "" {
		switch p.Type {
		case "":
			if p.To.Full == "" {
				p.broadcastPresence()
			}
		case "subscribe": p.subscribePresenceOut()
		case "subscribed": p.subscribedPresenceOut()
		}
	} else {
		switch p.Type {
		case "":
			p.directPresenceIn()
		case "subscribe": p.subscribePresenceIn()
		case "subscribed": p.subscribedPresenceIn()
		}
	}
}

func (p SMPacket) broadcastPresence() {
	session := sm.GetSession(p.Src.LocalResource())
	session.Available = true
	session.Presence = p.Stanza.Fragment // make hard copy? or pre-fragment and after- hang

	/* for _, item := range db.RosterItemsFrom(packet.From.Local) {
		stanza := &Stanza{
			Name: "presence",
			From: m.From,
			To: makeJid(item.Jid)}
		stanza.Fragment = m.Stanza.Fragment
		m.route(stanza)
	} */

	// TODO make it clearly
	p.From = p.Src
	p.Src = p.Src.BareJid()
	p.To = p.Src
	router.ch <- p.Packet

	// TODO send presence probes
}

func (p SMPacket) directPresenceIn() {
	// TODO handle set resource case
	user := p.To.Local

	for _, session := range sm.Sessions[user] {
		if !session.Available { continue }
		router.ch <- &Packet{Dest: session.Jid, Stanza: p.Stanza}
	}
}

func (p SMPacket) subscribePresenceOut() {
	state := db.GetSubState(p.To.Local, p.From.Bare)
	if state.IsOutNo() {
		state.SetOutPending()
		sm.pushRoster(p.To.Local, p.From.Bare)
	}

	p.Src = p.Src.BareJid()
	router.ch <- p.Packet
}

func (p SMPacket) subscribePresenceIn() {
	user, contact := p.To.Local, p.From.Bare
	state := db.GetSubState(user, contact)
	if state.IsInNo() {
		state.SetInPending()

		if sm.HasAvailable(user) {
			for _, session := range sm.Sessions[user] {
				if !session.Available { continue }
				router.ch <- &Packet{Dest: session.Jid, Stanza: p.Stanza}
			}
		} else {
			// TODO store presence
		}
	} else if state.IsInYes() {
		p.Swap()
		p.Type = "subscribed"
		p.Fragment = xml.NewBuilder().End()
		router.ch <- p.Packet
	}
}

func (p *SMPacket) subscribedPresenceOut() {
	user, contact := p.From.Local, p.To.Bare

	state := db.GetSubState(user, contact)
	if state.IsInPending() {
		state.SetInYes()

		router.ch <- &Packet{Src: p.Src.BareJid(), Dest: p.To, Stanza: p.Stanza}

		sm.pushRoster(user, contact)

		for _, session := range sm.Sessions[user] {
			if !session.Available { continue }
			// Id ? it is in rfc 6121
			stanza := &Stanza{Name: "presence", From: session.Jid, To: p.To, Fragment: session.Presence}
			router.ch <- &Packet{Src: session.Jid.BareJid(), Dest: p.To, Stanza: stanza}
		}
	}
}

func (p *SMPacket) subscribedPresenceIn() {
	user, contact := p.From.Local, p.To.Bare

	state := db.GetSubState(user, contact)
	if state.IsOutPending() {
		state.SetOutYes()

		for _, session := range sm.Sessions[user] {
			if !session.Interested { continue }
			router.ch <- &Packet{Dest: session.Jid, Stanza: p.Stanza}
		}

		sm.pushRoster(user, contact)
	}
}

func (sm *SM) pushRoster(user, contact string) {
	item := db.GetRosterItem(user, contact)

	builder := xml.NewBuilder().StartElement("query", "xmlns", rosterNs)
	item.WriteToBuilder(builder)
	fragment := builder.End()

	for _, session := range sm.Sessions[user] {
		if !session.Interested { continue }
		stanza := &Stanza{Name: "iq", Type: "set", Fragment: fragment}
		router.ch <- &Packet{Dest: session.Jid, Stanza: stanza}
	}
}
