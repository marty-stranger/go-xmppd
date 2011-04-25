package main

import (
	"crypto/rand"
	"fmt"
	"net"
	"os"

	"g/xml"
)

const (
	streamNs = "http://etherx.jabber.org/streams"
	streamErrorNs = "urn:ietf:params:xml:ns:xmpp-streams"
)

type Stream struct {
	net.Conn
	*xml.Reader
	*xml.Writer

	To	Jid
	Id	string
}

func newStream(netConn net.Conn) *Stream {
	s := &Stream{Conn: netConn}
	s.Reader = xml.NewReader(s)
	s.Writer = xml.NewWriter(s)
	return s
}


func (s *Stream) WriteStanza(stanza *Stanza) {
	s.StartElement(stanza.Kind.String(), "from", stanza.From.Full, "id", stanza.Id,
		"to", stanza.To.Full, "type", stanza.Type.String()).
		Raw(stanza.Fragment.String()).
		End()
}

func (s *Stream) Read(b []byte) (int, os.Error) {
	n, e := s.Conn.Read(b)
	debugln(s.To, string(b[:n]))
	return n, e
}

func (s *Stream) Write(b []byte) (int, os.Error) {
	n, e := s.Conn.Write(b)
	debugln(s.To, string(b))
	return n, e
}

func (s *Stream) writeStreamError(error string) {
	s.StartElement("stream:error").
		Element(error, "xmlns", streamErrorNs).
		End()
}

func (s *Stream) streamRecover() {
	error := recover()

	switch error := error.(type) {
	case nil:
		s.EndElementNoTrack("stream:stream")
		s.Close()
	case StreamError:
		s.writeStreamError(string(error))
		s.EndElementNoTrack("stream:stream")
		s.Close()
	case xml.ParserError:
		s.writeStreamError("not-well-formed")
		s.EndElementNoTrack("stream:stream")
		s.Close()
	// case xml.ReaderError:
	case net.Error: // or ?
		debugln(fmt.Sprintf("%#v", error))
	default:
		debugln(fmt.Sprintf("%#v", error))
		s.writeStreamError("internal-server-error")
		s.EndElementNoTrack("stream:stream")
		s.Close()
	}
}

type StreamError string

func (s *Stream) streamError(error string) {
	panic(StreamError(error))
}

func (s *Stream) generateStreamId() {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	s.Id = fmt.Sprintf("%x", bytes)
}
