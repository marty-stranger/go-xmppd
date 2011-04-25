package main

type C2SStream struct {
	*Stream
}

func (s *C2SStream) run() {
	defer s.streamRecover()

	s.negotiate()

	for {
		s.stanza()
	}
}
