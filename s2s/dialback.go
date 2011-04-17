package s2s

import (
	"hash"
	"crypto/sha256"
	"crypto/hmac"
	"fmt"

	"g/xmppd/base"
)

const dialbackSecret = "secret"

func dialbackKey(from, to, id string) string {
	var hash hash.Hash
	hash = sha256.New()
	hash.Write([]byte(dialbackSecret))
	key := fmt.Sprintf("%x", hash.Sum())

	message := fmt.Sprintf("%s %s %s", from, to, id)

	hash = hmac.NewSHA256([]byte(key))
	hash.Write([]byte(message))

	return fmt.Sprintf("%x", hash.Sum())
}

func (s *Stream) dialback() {
	key := dialbackKey(base.ServerName, s.streamTo, s.streamId)

	s.Element("db:result", "from", base.ServerName, "to", s.streamTo, key).End()

	cursor := s.ReadElement().Cursor()
	// should be db:result

	resultType := cursor.MustAttr("type")
	if resultType == "valid" {
		println("ok")
	}

	s.Element("presence", "from", "pavel@gxmppd.dyndns.org", "to", "pavel.rosputko@gmail.com",
		"type", "subscribe").End()

	s.ReadElement()
}

func (s *Stream) dialbackAccept() {
	// should be db result
	cursor := s.ReadElement().Cursor()

	from := cursor.MustAttr("from")
	to := cursor.MustAttr("to")

	// TODO verify @to, check hash
	if true {
		s.Element("db:result", "from", to, "to", from, "type", "valid").End()
	}

	// should be db verify
	cursor = s.ReadElement().Cursor()

	from = cursor.MustAttr("from")
	to = cursor.MustAttr("to")
	id := cursor.MustAttr("id")

	key := dialbackKey(to, from, id)

	var verifyType string
	if cursor.MustChars() == key {
		verifyType = "valid"
	} else {
		verifyType = "invalid"
	}

	s.Element("db:verify", "from", to, "id", id, "to", from, "type", verifyType).End()
}
