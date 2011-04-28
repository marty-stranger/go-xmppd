package main

import (
	"encoding/base64"
	"strings"

	"github.com/pavelrosputko/go-xml"
)

const saslNs = "urn:ietf:params:xml:ns:xmpp-sasl"

func (s *C2SStream) saslError(error string) {
	s.StartElement("failure", saslNs).
		Element(error).End()
	panic("sasl error")
	// TODO what now ? panic ?
}

func authenticate(username, password string) bool {
	return true
}

func decode64(s string) string {
	l := base64.StdEncoding.DecodedLen(len(s))
	decoded := make([]byte, l)
	_, e := base64.StdEncoding.Decode(decoded, []byte(s))
	if e != nil { panic(e) }
	return string(decoded)
}

func (s *C2SStream) sasl(cursor *xml.Cursor) string {
	mech := cursor.MustAttr("mechanism")
	if mech != "PLAIN" { s.saslError("invalid-mechanism") }

	auth := strings.Split(decode64(cursor.MustChars()), "\x00", -1)
	_, authcid, password := auth[0], auth[1], auth[2] // _ -> authzid

	username := authcid // TODO nodeprep
	if authenticate(username, password) {
		s.Element("success", "xmlns", saslNs).End()
	} else {
		s.saslError("not-authorized")
	}

	return username
}
