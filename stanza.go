package main

import (
	"fmt"

	"g/xml"
)

// NOTE all these consts will be useless if "abcd" == "abcd" would just compare strings' pointer and length

const (
	MessageKind = iota
	PresenceKind
	IQKind
)

const (
	ChatType = iota
	NoneType
	SubscribeType
	SubscribedType
	UnsubscribeType
	UnsubscribedType

	ProbeType
	GetType
	SetType
	ResultType

	ErrorType
)

type StanzaKind int

var StanzaKindsStrings = [...]string{MessageKind: "message", PresenceKind: "presence", IQKind: "iq"}

func (k StanzaKind) String() string {
	return StanzaKindsStrings[k]
}

type StanzaType int

var StanzaTypesStrings = [...]string{
	ChatType: "chat",

	NoneType: "",
	SubscribeType: "subscribe",
	SubscribedType: "subscribed",
	UnsubscribeType: "unsubscribe",
	UnsubscribedType: "unsubscribed",

	ProbeType: "probe",
	GetType: "get",
	SetType: "set",
	ResultType: "result",

	ErrorType: "error",
}

func (t StanzaType) String() string {
	return StanzaTypesStrings[t]
}

type Stanza struct {
	Kind	StanzaKind
	From	Jid
	Id	string
	To	Jid
	Type	StanzaType
	Lang	string
	*xml.Fragment
}

func (s *Stanza) String() string {
	return fmt.Sprintf("Name:%s From:%s Id:%s To:%s Type:%s %s",
		s.Kind, s.From.Full, s.Id, s.To.Full, s.Type, s.Fragment)
}


var TypesMap = map[StanzaKind]map[string]StanzaType{
	MessageKind: map[string]StanzaType{
		"chat": ChatType,
		"error": ErrorType,
	},
	PresenceKind: map[string]StanzaType{
		"": NoneType,
		"subscribe": SubscribeType,
		"subscribed": SubscribedType,
		"unsubscribe": UnsubscribeType,
		"unsubscribed": UnsubscribedType,
		"probe": ProbeType,
		"error": ErrorType,
	},
	IQKind: map[string]StanzaType{
		"get": GetType,
		"set": SetType,
		"result": ResultType,
		"error": ErrorType,
	},
}

// NOTE is switch ok? or use map? will it be faster? optimize via len comparisons
func newStanza(fragment *xml.Fragment) *Stanza {
	cursor := fragment.Cursor()

	stanza := &Stanza{}

	kind := cursor.Name()
	switch kind {
	case "message":
		stanza.Kind = MessageKind
	case "presence":
		stanza.Kind = PresenceKind
	case "iq":
		stanza.Kind = IQKind
	}

	stanzaTypeString, _ := cursor.Attr("type")
	stanzaType, foundStanzaType := TypesMap[stanza.Kind][stanzaTypeString]
	stanza.Type = stanzaType

	stanza.Id, _ = cursor.Attr("id")

	from, _ := cursor.Attr("from")
	to, _ := cursor.Attr("to")

	if from != "" { stanza.From = makeJid(from) }
	if to != "" { stanza.To = makeJid(to) }

	stanza.Fragment = cursor.ChildrenSlice()

	if !foundStanzaType {
		stanzaError(stanza, "modify", "bad-request")
	}

	return stanza
}

type StanzaError struct {
	*Stanza
	Type, Condition string
}

func stanzaError(stanza *Stanza, errorType, condition string) {
	panic(StanzaError{stanza, errorType, condition})
}

func (s *Stanza) Swap() {
	s.From, s.To = s.To, s.From
}
