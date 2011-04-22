package main

func (m *SM) presence(packet *Packet) {
	if packet.Src.Resource != "" {
		if packet.To.Full == "" && packet.Type == "" {
			m.broadcastPresence(packet)
		}
	} else {
		if packet.To.Full == "" && packet.Type == "" {
			m.directPresenceIn(packet)
		}
	}
}

func (m *SM) broadcastPresence(p *Packet) {
	// priority = 

	// presence := c.Stanza.Fragment.String()

	// presence := c.ChildrenString() // NOTE slicing does not let source be gced so copy of bytes or string()
	// fmt.Println("presence =", presence)

	// c.available = true

	/* for _, item := range db.RosterItemsFrom(packet.From.Local) {
		stanza := &Stanza{
			Name: "presence",
			From: m.From,
			To: makeJid(item.Jid)}
		stanza.Fragment = m.Stanza.Fragment
		m.route(stanza)
	} */

	resources := m.available[p.Src.Local]
	if resources == nil {
		resources = make(map[string]true)
		c.available[p.Src.Local] = resources
	}

	resources[p.Src.Resource] = true

	packet.Src = packet.Src.BareJid()
	router.ch <- packet

	// TODO send presence probes
}

func (m *SM) subscribePresenceOut() {
	m.route(m.stanza)

	state := db.SubState(m.To.Local, m.From.Bare)
	if state.isOutNo() {
		state.SetOutPending()

		item := db.RosterItem(m.To.Local, m.From.Bare)
		item.SetAsk()

		m.pushRoster(item, m.To.Local)
	}
}

func (m *SM) subscribePresenceIn(packet *Packet) {
	state := db.SubState(packet.To.Local, packet.From.Bare)
	if state.IsInNo() {
		state.SetInPending()

		if m.hasAvailable(packet.To.Local) {
			// TODO deliver to available
		} else {
			// TODO store presence
		}
	} else if state.IsInYes() {
		/* // auto-reply
		stanza := xml.New()
		stanza.Element("presence", "from", stanza.To().String(), "to", stanza.From().String(),
			"type", "subscribed")
		xmppd.Route(stanza) */
	}
}

func (m *SMOut) pushRoster(item *RosterItem, local string) {
	for resource := range m.interested(local) {
		stanza := Stanza{
			Name: "iq",
			To: makeJid(local, serverName, resource),
			Id: m.nextId(),
			Type: "set"}

		stanza.StartElement("query", "xmlns", rosterNs).
				Raw(item.xml()).
			EndElement()

		m.route(stanza)
	}
}

func (m *SM) directPresenceOut(packet *Packet) {
	state := db.SubState(packet.To.Local, packet.From.Bare)
	if !(state.IsInYes() && m.isAvailable(packet.From.Local, packet.From.Resource)) {
                // TODO we need to remember to send unavailable when client goes off!
		// PROPOSE c.directList append @t
	}

	packet.Src.unsetResource()
	routerCh <- packet
}

func (m *SM) directPresenceIn(packet *Packet) {
	for _, resource := range m.available(packet.To.Local, packet.To.Resource) {
	}
}

