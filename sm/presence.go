package main

func (p SMPacket) presence() {
	if p.Src.Domain == serverName && p.Src.Resource != "" {
		switch p.Type {
		case NoneType:
			if p.To.Full == "" {
				p.broadcastPresence()
			} else {
				p.directPresenceOut()
			}
		case SubscribeType: p.subscribePresenceOut()
		case SubscribedType: p.subscribedPresenceOut()
		case UnavailableType:
			if p.To.Full == "" {
				p.broadcastUnavailablePresence()
			}
		}
	} else {
		switch p.Type {
		case NoneType, UnavailableType: p.directPresenceIn()
		case SubscribeType: p.subscribePresenceIn()
		case SubscribedType: p.subscribedPresenceIn()
		case ProbeType: p.probePresenceIn()
		}
	}
}

func (p SMPacket) broadcastPresence() {
	session := sm.GetSession(p.Src.LocalResource())
	session.Presence = p.Stanza.Fragment // make hard copy? or pre-fragment and after- hang

	bareSrc := p.Src.BareJid()
	user := p.Src.Local

	for _, contact := range db.SubStatesInYes(user) {
		jid := makeJid(contact)
		stanza := &Stanza{Kind: PresenceKind, From: p.Src, To: jid, Fragment: session.Presence}
		router.Ch <- &Packet{bareSrc, jid, stanza}
	}

	// NOTE join with previous for ? // send only what we do not know about
	// TODO send presence probes 
	// TODO how about contact resides this server ?
	if !session.Available {
		session.Available = true
		for _, contact := range db.SubStatesOutYes(user) {
			jid := makeJid(contact)
			stanza := &Stanza{Kind: PresenceKind, From: bareSrc, Id: sm.nextId(), To:jid,
				Type: ProbeType, Fragment: nil}
			router.Ch <- &Packet{bareSrc, jid, stanza}
		}
	}

	p.From, p.Src, p.To = p.Src, p.Src.BareJid(), p.Src
	router.Ch <- p.Packet
}

func (p SMPacket) broadcastUnavailablePresence() {
	// TODO verify elements ? MUST NOT has show, priority
	user := p.Src.Local
	bareSrc := p.Src.BareJid()

	// double unavailable ? do what ?

	for _, contact := range db.SubStatesInYes(user) {
		jid := makeJid(contact)
		stanza := &Stanza{Kind: PresenceKind, From: p.Src, To: jid, Type: UnavailableType,
			Fragment: p.Fragment}
		router.Ch <- &Packet{bareSrc, jid, stanza}
	}

	session := sm.GetSession(p.Src.LocalResource())
	session.Unavailable()

	p.From, p.Src, p.To = p.Src, p.Src.BareJid(), p.Src
	router.Ch <- p.Packet
}

func (p SMPacket) directPresenceOut() {
	// policy to use just From, To ?
	user, contact := p.From.Local, p.To.Bare
	state := db.GetSubState(user, contact)

	session := sm.GetSession(p.Src.LocalResource())

	if !(state.IsInYes() && session.Available) {
		// TODO maintain direct list, modify broadcastUnavailable, modify probePresence
	}

	p.Src = p.Src.BareJid()
	router.Ch <- p.Packet
}

func (p SMPacket) directPresenceIn() {
	// TODO handle non empty resorces case 

	user := p.To.Local

	for _, session := range sm.Sessions[user] {
		if !session.Available { continue }
		router.Ch <- &Packet{Dest: session.Jid, Stanza: p.Stanza}
	}
}

func (p SMPacket) probePresenceIn() {
	// NOTE from should/must be bare or not ?
	user, contact := p.To.Local, p.From.Bare

	// NOTE handling of non existent user case ?
	state := db.GetSubState(user, contact)
	if !state.IsInYes() {
		// assure p.To is Bare or presence leak
		p.Swap()
		p.Type, p.Fragment = UnsubscribedType, nil
		router.Ch <- p.Packet
		return
	}

	// TODO temporary move
	if p.To.Resource == "" {
		if !sm.HasAvailable(user) {
			p.Swap()
			p.Type, p.Fragment = UnavailableType, nil
			router.Ch <- p.Packet
			return
		}

		for _, session := range sm.Sessions[user] {
			if !session.Available { continue }

			jid := session.Jid
			stanza := &Stanza{
				Kind: PresenceKind,
				From: jid,
				// Id: sm.nextId(), // http://tools.ietf.org/html/rfc6121#section-4.3.2.1 says
						// to send id of initial presence ? store it in session ?
				To: p.From,
				Fragment: session.Presence,
			}

			router.Ch <- &Packet{jid.BareJid(), p.From, stanza}
		}
	} else {
		session := sm.GetSession(p.To.LocalResource())
		if session != nil {
			stanza := &Stanza{From: session.Jid, To: p.From}
			router.Ch <- &Packet{session.Jid, p.From, stanza}
		} else {
			// TODO what ?
		}
	}
}
