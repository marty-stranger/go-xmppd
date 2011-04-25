package main

import (
	"g/xml"
)

func (d *Db) WriteRoster(user string, builder *xml.Builder) {
	substatesValues := d.Hgetall("substates:" + user)

	substates := make(map[string]SubState)
	for i := 0; i < len(substatesValues); i += 2 {
		jid := string(substatesValues[i])
		debugln("jid =", jid)
		substate := substatesValues[i + 1][0]
		substates[jid] = SubState(substate)
		debugln("substate =", SubState(substate))
	}

	itemsValues := d.Hgetall("roster.items:" + user)
	// groupsValues := d.Smembers("roster.groups:" + user)

	println("len(itemsValues) =", len(itemsValues))
	for i := 0; i < len(itemsValues); i += 2 {
		jid := string(itemsValues[i])
		debugln("jid =", jid)
		name := string(itemsValues[i + 1])
		debugln("name =", name)

		substate := substates[jid]
		subscription, ask := substate.SubscriptionAsk()

		builder.Element("item",
			"ask", ask,
			"jid", jid,
			"name", name,
			"subscription", subscription)
		// TODO groups
	}
}

type RosterItemData struct {
	Name, Subscription, Ask string
	Groups []string
}

type RosterItem struct {
	User, Jid string
	RosterItemData
}

func (d *Db) DeleteRosterItem(user, contact string) {
	d.Hdel("roster.items:" + user, contact)
	d.Hdel("substates:" + user, contact)
	// TODO groups
}

func (d *Db) GetRosterItem(user, contact string) *RosterItem {
	itemValue := d.Hget("roster.items:" + user, contact)
	substateValue := d.Hget("substates:" + user, contact)

	item := RosterItem{User: user, Jid: contact}

	if itemValue != nil {
		item.Name = string(itemValue)
	}

	var substate SubState
	if substateValue != nil {
		substate = SubState(substateValue[0])
	}

	item.Subscription, item.Ask = substate.SubscriptionAsk()
	// TODO groups

	return &item
}

func (i *RosterItem) UpdateFromFragment(fragment *xml.Fragment) {
	cursor := fragment.Cursor()
	name, _ := cursor.Attr("name")
	// TODO groups

	i.Name = name
	db.Hset("roster.items:" + i.User, i.Jid, name)
	db.Hset("substates:" + i.User, i.Jid, string([]byte{0}))
}

func (i *RosterItem) WriteToBuilder(builder *xml.Builder) {
	builder.Element("item",
		"ask", i.Ask,
		"jid", i.Jid,
		"name", i.Name,
		"subscription", i.Subscription)
}
