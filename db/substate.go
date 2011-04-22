package main

type SubState uint8

func (s SubState) IsInNo() bool { return s & 4 == 0 }
func (s SubState) IsInYes() bool { return s & 4 == 1 }
func (s SubState) IsInPending() bool { return s & 4 == 2 }

func (s SubState) IsOutNo() bool { return s >> 2 & 4 == 0 }
func (s SubState) IsOutYes() bool { return s >> 2 & 4 == 1 }
func (s SubState) IsOutPending() bool { return s >> 2 & 4 == 2 }

var subscriptionsNames = []string{"none", "from", "to", "both"}

func (s SubState) SubscriptionAsk() (subscription string, ask string) {
	in, out := s & 4, s >> 2 & 4

	if out == 2 { ask = "subscribe" }

	if in == 2 { in = 0 }
	if out == 2 { out = 0 }

	subscription = subscriptionsNames[in + 2*out]
	return
}

type SubStateDbItem struct {
	user, contact string
	SubState
}

func (db *Db) GetSubState(user, contact string) SubStateDbItem {
	bytes := db.Hget("substates:" + user, contact)

	var value uint8
	if bytes != nil { value = bytes[0] }

	return SubStateDbItem{user, contact, SubState(value)}
}

func (s *SubStateDbItem) SetInNo() { s.SubState &^= 1; s.SubState &^= 2; s.Save() }
func (s *SubStateDbItem) SetInYes() { s.SubState |= 1; s.SubState &^= 2; s.Save() }
func (s *SubStateDbItem) SetInPending() { s.SubState &^= 1; s.SubState |= 2; s.Save() }

func (s *SubStateDbItem) SetOutNo() { s.SubState &^= 4; s.SubState &^= 8; s.Save() }
func (s *SubStateDbItem) SetOutYes() { s.SubState |= 4; s.SubState &^= 8; s.Save() }
func (s *SubStateDbItem) SetOutPending() { s.SubState &^= 4; s.SubState |= 8; s.Save() }

func (s *SubStateDbItem) Save() {
	db.Hset("s:" + s.user, s.contact, string([]byte{uint8(s.SubState)}))
}
