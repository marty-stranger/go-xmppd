package main

import (
	"github.com/pavelrosputko/go-xml"
)

const rosterNs = "jabber:iq:roster"

func (p SMPacket) rosterIQ() {
	debugln("")
	switch p.Type {
	case GetType:
		sm.GetSession(p.Src.LocalResource()).Interested = true

		builder := xml.NewBuilder().
			StartElement("query", "xmlns", rosterNs)
		db.WriteRoster(p.Src.Local, builder)
		fragment := builder.End()

		p.Swap()
		p.Type = ResultType
		p.Fragment = fragment
		router.Ch <- p.Packet
	case SetType:
		debugln("")
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
			stanza := &Stanza{Kind: IQKind, Id: sm.nextId(), Type: SetType, Fragment: fragment}
			packet := &Packet{Dest: session.Jid, Stanza: stanza}
			router.Ch <- packet
		}

		p.Swap()
		p.Type = ResultType
		p.Fragment = xml.NewBuilder().End()
		router.Ch <- p.Packet
	}
}
