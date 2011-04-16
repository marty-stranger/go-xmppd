package main

import (
	"strings"
)

type Jid struct {
	Full, Bare, Local, Domain, Resource string
}

func makeJid(s string) Jid {
	// TODO nodeprep

	j := Jid{}

	j.Full = s

	slashIndex := strings.Index(s, "/")
	j.Bare, j.Resource = s[:slashIndex], s[slashIndex + 1:]

	atIndex := strings.Index(j.Bare, "@")
	j.Local, j.Domain = j.Bare[:atIndex], j.Bare[atIndex + 1:]

	return j
}
