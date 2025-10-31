package networking

import (
	"context"
	"log/slog"
	"net"
)

type SessionState interface {
	Handle(ctx context.Context) error
}

type ClientInfo struct {
	playerName string
}

type Session struct {
	ID     uint64
	Logger *slog.Logger

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
