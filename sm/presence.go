package main

func (m *SM) presence(packet *Packet) {
	if packet.Src.Resource != "" {
		if packet.To.Full == "" && packet.Type == "" {
			m.broadcastPresence(packet)
		}
	} else {
		println("sm presence in")
		if packet.Type == "" {
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
		resources = make(map[string]bool)
		m.available[p.Src.Local] = resources
	}

	resources[p.Src.Resource] = true

	// TODO make it clearly
	p.From = p.Src
	p.Src = p.Src.BareJid()
	p.To = p.Src
	router.ch <- p

	// TODO send presence probes
}

func (m *SM) directPresenceIn(p *Packet) {
	// TODO handle set resource case
	resources := m.available[p.To.Local]
	if resources == nil { return }

	for resource, _ := range resources {
		pp := &Packet{}
		pp.Src = p.Src
		pp.Dest = makeJid(p.To.Local + "@" + serverName + "/" + resource)
		pp.Stanza = p.Stanza
		router.ch <- pp
	}
}

