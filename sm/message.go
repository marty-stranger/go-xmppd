package main

func (m *SM) message(packet *Packet) {
	if packet.Src.Resource != "" {
		packet.Src.EmptyResource()
		routerCh <- packet
		return
	}

	if m.Stanza.To.Resource == "" {
		m.bareMessage(packet)
	} else {
		m.fullMessage(packet)
	}
}

func (m *SM) bareMessage(packet *Packet)
	if m.hasConnected(packet.To.Local) {
		switch packet.Type {
		case "normal", "chat" {
			nonNegative := m.getNonNegative(packet.To.Local)
			switch len(nonNegative) {
			case 0:
				// TODO store offline
			case 1:
				resource := nonNegative[0]
				packet.To.SetResource(resource)
				routerCh <- packet
			default:
				// TODO choose one and deliver to it or to all
			}
		case "groupchat":
		case "headline":
		case "error":
		}
	} else {
		switch packet.Type {
		case "normal", "chat": // TODO offline
		case "groupchat": // TODO service-unavailable
		case "headline", "error": // ignore
		}
	}
}

func (m *SM) fullMessage(packet *Packet)
	if m.hasConnected(packet.To.Local, packet.To.Resource) {
		m.route()
	} else {
		switch packet.Type {
		case "normal", "groupchat", "headline": // TODO offline
		case "chat": m.bareMessage(packet)
		case "error": // ignore
	}
}
