package main

import "fmt"

type SubState uint8

func (s SubState) IsInNo() bool { return s & 3 == 0 }
func (s SubState) IsInYes() bool { return s & 3 == 1 }
func (s SubState) IsInPending() bool { return s & 3 == 2 }

func (s SubState) IsOutNo() bool { return s >> 2 & 3 == 0 }
func (s SubState) IsOutYes() bool { return s >> 2 & 3 == 1 }
func (s SubState) IsOutPending() bool { return s >> 2 & 3 == 2 }

var subscriptionsNames = []string{"none", "from", "to", "both"}

func (s SubState) SubscriptionAsk() (subscription string, ask string) {
	in, out := s & 3, s >> 2 & 3
	debugln(s)

	if out == 2 { ask = "subscribe" }

	if in == 2 { in = 0 }
	if out == 2 { out = 0 }

	subscription = subscriptionsNames[in + 2*out]
	return
}

var substatesNames = []string{"no", "yes", "pending"}

func (s SubState) String() string {
	in, out := s & 3, s >> 2 & 3
	return fmt.Sprintf("in:%s out:%s", substatesNames[in], substatesNames[out])
}

type SubStateDbItem struct {
	user, contact string
	SubState
}

func (db *Db) GetSubState(user, contact string) SubStateDbItem {
	debugln(user, contact)
	bytes := db.Hget("substates:" + user, contact)

	var value uint8
	if bytes != nil { value = bytes[0] }

	debugln(value)

	return SubStateDbItem{user, contact, SubState(value)}
}

func (s *SubStateDbItem) SetInNo() { s.SubState &^= 1; s.SubState &^= 2; s.Save() }
func (s *SubStateDbItem) SetInYes() { s.SubState |= 1; s.SubState &^= 2; s.Save() }
func (s *SubStateDbItem) SetInPending() { s.SubState &^= 1; s.SubState |= 2; s.Save() }

func (s *SubStateDbItem) SetOutNo() { s.SubState &^= 4; s.SubState &^= 8; s.Save() }
func (s *SubStateDbItem) SetOutYes() { s.SubState |= 4; s.SubState &^= 8; s.Save() }
func (s *SubStateDbItem) SetOutPending() { s.SubState &^= 4; s.SubState |= 8; s.Save() }

func (s *SubStateDbItem) Save() {
	debugln(s)
	db.Hset("substates:" + s.user, s.contact, string([]byte{uint8(s.SubState)}))
}

func (d *Db) SubStatesInYes(user string) []string {
	substatesValues := d.Hgetall("substates:" + user)

	jids := make([]string, 0, len(substatesValues))
	for i := 0; i < len(substatesValues); i += 2 {
		jid := string(substatesValues[i])
		substate := SubState(substatesValues[i + 1][0])

		if substate.IsInYes() {
			jids = append(jids, jid)
		}
	}

	return jids
}

func (d *Db) SubStatesOutYes(user string) []string {
	substatesValues := d.Hgetall("substates:" + user)

	jids := make([]string, 0, len(substatesValues))
	for i := 0; i < len(substatesValues); i += 2 {
		jid := string(substatesValues[i])
		substate := SubState(substatesValues[i + 1][0])

		if substate.IsOutYes() {
			jids = append(jids, jid)
		}
	}

	return jids
}
