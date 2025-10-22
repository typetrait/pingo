package networking

import (
	"net"
)

type SessionState interface {
	Handle() error
}

type ClientInfo struct {
}

type Session struct {
	server     *Server
	clientInfo *ClientInfo
	conn       net.Conn
	state      SessionState
}

func NewSession(server *Server, clientInfo *ClientInfo, conn net.Conn) *Session {
	s := &Session{
		server:     server,
		clientInfo: clientInfo,
		conn:       conn,
	}
	s.state = &NegotiateSessionState{s}
	return s
}

func (s *Session) ClientInfo() *ClientInfo {
	return s.clientInfo
}

func (s *Session) SetState(state SessionState) {
	s.state = state
}
