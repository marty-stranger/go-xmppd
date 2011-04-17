package main

const bindNs = "urn:ietf:params:xml:ns:xmpp-bind"

func (c *C2SConn) bind(local string) {
	id := c.MustAttr("id")
	c.MustToChild()

	var resource string
	if c.ToChild() { resource = c.MustChars() }

	if resource == "" {
		// TODO generate
		resource = "gxmppd"
	}

	// TODO make jid
	jid := local + "@" + serverName + "/" + resource
	c.jid = makeJid(jid)

	// TODO check resource availabily
	c.StartElement("iq", "id", id, "type", "result").
		StartElement("bind", "xmlns", bindNs).
			Element("jid", jid).End()

	// TODO add to c2sConns
	c2sConns.Add(c.jid.Local, c.jid.Resource, c)

	c.connected = true
}

