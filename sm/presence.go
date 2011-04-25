package main

import (
	"g/xml"
)

func (p SMPacket) presence() {
	if p.Src.Domain == serverName && p.Src.Resource != "" {
		switch p.Type {
		case NoneType:
			if p.To.Full == "" {
				p.broadcastPresence()
			}
		case SubscribeType: p.subscribePresenceOut()
		case SubscribedType: p.subscribedPresenceOut()
		}
	} else {
		switch p.Type {
		case NoneType:
			p.directPresenceIn()
		case SubscribeType: p.subscribePresenceIn()
		case SubscribedType: p.subscribedPresenceIn()
		}
	}
}

func (p SMPacket) broadcastPresence() {
	session := sm.GetSession(p.Src.LocalResource())
	session.Available = true
	session.Presence = p.Stanza.Fragment // make hard copy? or pre-fragment and after- hang

	for _, jid := range db.SubStatesInYes(p.From.Local) {
		stanza := &Stanza{
			Kind: PresenceKind,
			From: p.Src,
			To: makeJid(jid), // NOTE no prepare is needed
			Fragment: session.Presence,
		}

		router.Ch <- &Packet{Src: p.Src.BareJid(), Dest: stanza.To, Stanza: stanza}
	}

	// TODO make it clearly
	p.From = p.Src
	p.Src = p.Src.BareJid()
	p.To = p.Src
	router.Ch <- p.Packet

	// TODO send presence probes, TODO send only to contact with do not know about
	// TODO how about contact resides with server ?
	for _, jid := range db.SubStatesOutYes(p.From.Local) {
		stanza := &Stanza{
			Kind: PresenceKind,
			From: p.Src.BareJid(),
			Id: sm.nextId(),
			To: makeJid(jid), // NOTE no prepare is needed
			Type: ProbeType,
			Fragment: xml.NewBuilder().End(),
		}

		router.Ch <- &Packet{Src: p.Src.BareJid(), Dest: stanza.To, Stanza: stanza}
	}
}

func (p SMPacket) directPresenceIn() {
	// TODO handle set resource case
	user := p.To.Local

	for _, session := range sm.Sessions[user] {
		if !session.Available { continue }
		router.Ch <- &Packet{Dest: session.Jid, Stanza: p.Stanza}
	}
}

func (p SMPacket) subscribePresenceOut() {
	debugln("")
	user, contact := p.From.Local, p.To.Bare
	state := db.GetSubState(user, contact)
	debugln(state)
	if state.IsOutNo() {
		debugln("")
		state.SetOutPending()
		sm.pushRoster(user, contact)
	}

	p.Dest = p.To
	router.Ch <- p.Packet
}

func (p SMPacket) subscribePresenceIn() {
	user, contact := p.To.Local, p.From.Bare
	state := db.GetSubState(user, contact)
	if state.IsInNo() {
		state.SetInPending()

		if sm.HasAvailable(user) {
			for _, session := range sm.Sessions[user] {
				if !session.Available { continue }
				router.Ch <- &Packet{Dest: session.Jid, Stanza: p.Stanza}
			}
		} else {
			// TODO store presence
		}
	} else if state.IsInYes() {
		p.Swap()
		p.Type = SubscribedType
		p.Fragment = xml.NewBuilder().End()
		router.Ch <- p.Packet
	}
}

func (p *SMPacket) subscribedPresenceOut() {
	user, contact := p.From.Local, p.To.Bare

	state := db.GetSubState(user, contact)
	if state.IsInPending() {
		state.SetInYes()

		router.Ch <- &Packet{Src: p.Src.BareJid(), Dest: p.To, Stanza: p.Stanza}

		sm.pushRoster(user, contact)

		for _, session := range sm.Sessions[user] {
			if !session.Available { continue }
			// Id ? it is in rfc 6121
			stanza := &Stanza{Kind: PresenceKind, From: session.Jid, To: p.To, Fragment: session.Presence}
			router.Ch <- &Packet{Src: session.Jid.BareJid(), Dest: p.To, Stanza: stanza}
		}
	}
}

func (p *SMPacket) subscribedPresenceIn() {
	user, contact := p.To.Local, p.From.Bare

	state := db.GetSubState(user, contact)
	if state.IsOutPending() {
		state.SetOutYes()

		for _, session := range sm.Sessions[user] {
			if !session.Interested { continue }
			router.Ch <- &Packet{Dest: session.Jid, Stanza: p.Stanza}
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
		stanza := &Stanza{Kind: IQKind, Type: SetType, Fragment: fragment}
		router.Ch <- &Packet{Dest: session.Jid, Stanza: stanza}
	}
}
