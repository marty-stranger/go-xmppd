package main

type C2SStream struct {
	*Stream
}

func (s *C2SStream) run() {
	defer s.streamRecover()

	defer func() {
		if s.To.Full != "" {
			sm.UnbindResource(s.To.LocalResource())
		}
	}()

	s.negotiate()

	for {
		s.stanza()
	}
}
