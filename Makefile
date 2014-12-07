include $(GOROOT)/src/Make.inc

TARG = xmppd

GOFILES = \
	logger.go \
	config.go \
	stream.go \
	jid.go \
	main.go \
	stanza.go \
	packet.go \
	router.go \
	c2s/c2s.go \
	c2s/stream.go \
	c2s/negotiate.go \
	c2s/tls.go \
	c2s/sasl.go \
	c2s/bind.go \
	c2s/stanza.go \
	s2s/s2s.go \
	s2s/stream.go \
	s2s/dialback.go \
	s2s/accept.go \
	sm/sm.go \
	sm/iq.go \
	sm/roster.go \
	sm/presence.go \
	sm/presence-subscr.go \
	local/local.go \
	db/db.go \
	db/roster.go \
	db/substate.go

include $(GOROOT)/src/Make.cmd
