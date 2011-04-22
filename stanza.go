package main

import (
	"fmt"

	"g/xml"
)

type Stanza struct {
	Name, Id, Type	string
	From, To	Jid

	*xml.Fragment
}

func newStanza(fragment *xml.Fragment) *Stanza {
	println("newStanza", fragment.String())
	cursor := fragment.Cursor()

	stanza := Stanza{}
	stanza.Name = cursor.Name()
	stanza.Id, _ = cursor.Attr("id")
	stanza.Type, _ = cursor.Attr("type")

	from, _ := cursor.Attr("from")
	to, _ := cursor.Attr("to")

	if from != "" { stanza.From = makeJid(from) }
	if to != "" { stanza.To = makeJid(to) }

	stanza.Fragment = cursor.ChildrenSlice()

	return &stanza
}

func (s *Stanza) String() string {
	return fmt.Sprintf("%s %s %s %s %s %s", s.Name, s.From.Full, s.Id, s.To.Full, s.Type, s.Fragment)
}
