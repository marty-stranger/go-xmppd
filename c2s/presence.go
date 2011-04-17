package main

import "fmt"

func (c *C2SConn) presence() {
	to, _ := c.Attr("to")
	presenceType, _ := c.Attr("type")

	if to == "" {
		if presenceType == "" {
			c.broadcastPresence()
		}
	}
}

func (c *C2SConn) broadcastPresence() {
	// c.priority = 

	presence := c.ChildrenString() // NOTE slicing does not let source be gced so copy of bytes or string()
	fmt.Println("presence =", presence)

	c.available = true

	// TODO route ? process as direct presence
	for _, cc := range c2sConns.Available(c.jid.Local) {
		fmt.Println("cc.jid", cc.jid)
		cc.StartElement("presence", "from", c.jid.Full, "to", cc.jid.Bare).
			Raw(presence).
			End()
	}

}
