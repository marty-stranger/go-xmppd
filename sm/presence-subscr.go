package main

import (
	"github.com/pavelrosputko/go-xml"
)

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

	// FIXME
	p.Dest = p.To
	router.Ch <- p.Packet
}

func (p SMPacket) subscribePresenceIn() {
	user, contact := p.To.Local, p.From.Bare
	state := db.GetSubState(user, contact)
	if state.IsInNo() {
		state.SetInPending()

		if sm.HasAvailable(user) {
			// propose sm.deliverAvailable(p.Stanza)
			for _, session := range sm.Sessions[user] {
				if !session.Available { continue }
				router.Ch <- &Packet{Dest: session.Jid, Stanza: p.Stanza}
			}
		} else {
			// TODO store presence
		}
	} else if state.IsInYes() {
		p.Swap()
		p.Type, p.Fragment = SubscribedType, nil
		router.Ch <- p.Packet
	}
}

func (p SMPacket) subscribedPresenceOut() {
	user, contact := p.From.Local, p.To.Bare

	state := db.GetSubState(user, contact)
	if state.IsInPending() {
		state.SetInYes()

		to := p.To
		router.Ch <- &Packet{Src: p.Src.BareJid(), Dest: to, Stanza: p.Stanza}

		sm.pushRoster(user, contact)

		for _, session := range sm.Sessions[user] {
			if !session.Available { continue }
			stanza := &Stanza{From: session.Jid, To: to, Fragment: session.Presence}
			router.Ch <- &Packet{session.Jid.BareJid(), to, stanza}
		}
	}
}

func (p SMPacket) subscribedPresenceIn() {
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

func (p SMPacket) unsubscirbePresenceOut() {
	user, contact := p.From.Local, p.To.Bare
	state := db.GetSubState(user, contact)

	p.Src = p.Src.BareJid()
	router.Ch <- p.Packet

	if !state.IsOutNo() {
		state.SetOutNo()
		sm.pushRoster(user, contact)
	}
}

func (p SMPacket) unsubscribePresenceIn() {
	user, contact := p.To.Local, p.From.Bare
	state := db.GetSubState(user, contact)

	if !state.IsInNo() {
		state.SetInNo()

		for _, session := range sm.Sessions[user] {
			if !session.Interested { continue }
			router.Ch <- &Packet{Dest: session.Jid, Stanza: p.Stanza}
		}

		sm.pushRoster(user, contact)
	}
}

func (p SMPacket) unsubscribedPresenceOut() {
	user, contact := p.From.Local, p.To.Bare

	state := db.GetSubState(user, contact)

	if state.IsInNo() {
		if state.PreAppr() {
			state.RemovePreAppr()
		}
	} else {
		from, to := p.From, p.To
		for _, session := range sm.Sessions[user] {
			if !session.Available { continue }
			// NOTE what id? generated ?
			stanza := &Stanza{Kind: PresenceKind, From: session.Jid, Id: sm.nextId(),
				To: to, Type: UnavailableType}
			router.Ch <- &Packet{from, to, stanza}
		}

		p.From, p.Src, p.Dest = p.Src.BareJid(), p.Src.BareJid(), p.From
		router.Ch <- p.Packet

		sm.pushRoster(user, contact)
	}
}

func (p SMPacket) unsubscribedPresenceIn() {
	user, contact := p.To.Local, p.From.Bare

	state := db.GetSubState(user, contact)
	if !state.IsOutNo() {
		state.SetOutNo()

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
