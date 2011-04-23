package main

import (
	"g/xml"
)

const rosterNs = "jabber:iq:roster"

func (p SMPacket) rosterIQ() {
	switch p.Type {
	case "get":
		sm.GetSession(p.Src.LocalResource()).Interested = true

		builder := xml.NewBuilder().
			StartElement("query", "xmlns", rosterNs)
		db.WriteRoster(p.Src.Local, builder)
		fragment := builder.End()

		p.Swap()
		p.Type = "result"
		p.Fragment = fragment
		router.ch <- p.Packet
	case "set":
		cursor := p.Cursor()
		cursor.MustToChild()

		jid := cursor.MustAttr("jid") // TODO prepare jid ?
		subscr, _ := cursor.Attr("subscription")

		builder := xml.NewBuilder().
			Element("query", "xmlns", rosterNs)

		user := p.Src.Local
		if subscr == "remove" {
			// TODO if no key return item-not-found
			db.DeleteRosterItem(user, jid)
			builder.Element("item", "jid", jid, "subscription", "remove")
		} else {
			item := db.GetRosterItem(user, jid)

			cursor = p.Cursor()
			item.UpdateFromFragment(cursor.ChildrenSlice())
			item.WriteToBuilder(builder)
		}

		fragment := builder.End()

		for _, session := range sm.Sessions[user] {
			if !session.Interested { continue }
			stanza := &Stanza{Name: "iq", Id: sm.nextId(), Type: "set", Fragment: fragment}
			packet := &Packet{Dest: session.Jid, Stanza: stanza}
			router.ch <- packet
		}

		p.Swap()
		p.Type = "result"
		p.Fragment = xml.NewBuilder().End()
		router.ch <- p.Packet
	}
}
