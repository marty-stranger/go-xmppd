package main

import (
	"strings"
)

type Jid struct {
	Full, Bare, Local, Domain, Resource string
}

func makeJid(s string) Jid {
	// TODO nodeprep
	println("makeJid: s =", s)

	j := Jid{}

	j.Full = s

	slashIndex := strings.Index(s, "/")
	if slashIndex != -1 {
		j.Bare, j.Resource = s[:slashIndex], s[slashIndex + 1:]
	} else {
		j.Bare = s
	}

	atIndex := strings.Index(j.Bare, "@")
	if atIndex != -1 {
		j.Local, j.Domain = j.Bare[:atIndex], j.Bare[atIndex + 1:]
	} else {
		j.Domain = j.Bare
	}

	return j
}

func (j Jid) BareJid() Jid {
	return Jid{
		Full: j.Bare,
		Bare: j.Bare,
		Local: j.Local,
		Domain: j.Domain,
		Resource: ""}
}
