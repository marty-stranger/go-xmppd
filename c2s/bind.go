package main

import (
	"g/xml"
)

const bindNs = "urn:ietf:params:xml:ns:xmpp-bind"

func (s *C2SStream) bind(local string, cursor *xml.Cursor) {
	id := cursor.MustAttr("id")

	cursor.MustToChild()

	var resource string
	if cursor.ToChild() {
		resource = cursor.MustChars()
	}

	if resource == "" {
		// TODO generate
		resource = "gxmppd"
	}

	// TODO think about rpc or stream for such communication
	if sm.BindResource(local, resource) {
		c2s.Add(local, resource, s)
		println("c2s, bind", local, resource)

		jid := local + "@" + serverName + "/" + resource
		s.To = makeJid(jid)
		s.StartElement("iq", "id", id, "type", "result").
			StartElement("bind", "xmlns", bindNs).
				Element("jid", jid).
			End()
	} else {
		// TODO
	}

}
