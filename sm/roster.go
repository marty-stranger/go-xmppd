package main

import (
	"fmt"

	"g/xml"
)

const rosterNs = "jabber:iq:roster"

func (m *SM) rosterIQ(packet *Packet) {
	switch packet.Type {
	case "get":
		m.interested.Add(packet.Src.LocalResource())
		fmt.Println("m.interested = ", m.interested)

		builder := xml.NewBuilder().
			StartElement("query", "xmlns", rosterNs)
		db.WriteRoster(packet.Src.Local, builder)
		fragment := builder.End()

		packet.Swap()
		packet.Type = "result"
		packet.Fragment = fragment
		router.ch <- packet
	case "set":
		println("roster-set")
		cursor := packet.Cursor()
		cursor.MustToChild()

		jid := cursor.MustAttr("jid") // TODO prepare jid ?
		subscr, _ := cursor.Attr("subscription")

		builder := xml.NewBuilder().
			Element("query", "xmlns", rosterNs)

		user := packet.Src.Local
		if subscr == "remove" {
			// TODO if no key return item-not-found
			db.DeleteRosterItem(user, jid)
			builder.Element("item", "jid", jid, "subscription", "remove")
		} else {
			item := db.GetRosterItem(user, jid)

			cursor = packet.Cursor()
			item.UpdateFromFragment(cursor.ChildrenSlice())
			item.WriteToBuilder(builder)
		}

		if m.interested[user] != nil {
			for resource, _ := range m.interested[user] {
				stanza := &Stanza{
					Name: "iq",
					To: makeJid(user + "@" + serverName + "/" + resource),
					Id: m.nextId(),
					Type: "set"}
				stanza.Fragment = builder.End()

				packet := &Packet{}
				packet.Dest = stanza.To
				packet.Stanza = stanza
				router.ch <- packet
			}
		}

		packet.Swap()
		packet.Type = "result"
		packet.Fragment = xml.NewBuilder().End()
		router.ch <- packet
	}
}
