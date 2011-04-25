package main

import "testing"

func TestSubState(t *testing.T) {
	var s SubStateDbItem

	s.SetInYes()
	if !s.IsInYes() { t.Fatal("") }

	s.SetOutPending()
	if !s.IsOutPending() { t.Fatal("") }
}
