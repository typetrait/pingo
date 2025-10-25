package networking

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/typetrait/pingo/cmd/server/game"
	"github.com/typetrait/pingo/internal/networking"
	"github.com/typetrait/pingo/internal/packet/serverbound"
)

const (
	protocolVersion uint8 = 0
)

var (
	ErrUnknownPacketType = errors.New("unknown packet type")
)

type Server struct {
	running bool

	mu       sync.Mutex
	sessions []*Session

	mu2     sync.Mutex
	matches map[string]*game.Match
}

func NewServer() *Server {
	return &Server{
		running:  false,
		sessions: make([]*Session, 0),
		matches:  make(map[string]*game.Match),
	}
}

func (s *Server) Start() error {
	addr := ":7777"

	slog.Info("listening", "address", addr)
	slog.Info("waiting for connections")
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("could not listen on address %s", addr)
	}

	s.running = true
	for s.running {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("could not accept connection")
			continue
		}

		go s.HandleConnection(conn)
	}

	return nil
}

func (s *Server) Stop() {
	s.running = false
}

func (s *Server) HandleConnection(conn net.Conn) {
	slog.Debug("handling new connection")

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			slog.Error("could not close connection")
		}
	}(conn)

	session := NewSession(s, &ClientInfo{}, conn)
	sessionID := s.addSession(session)

	session.Logger = slog.
		Default().
		WithGroup(fmt.Sprintf("session_%s", strconv.Itoa(sessionID)))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		session.Logger.Info("session state", "state", session.state)
		err := session.state.Handle(ctx)
		if err != nil {
			if errors.Is(err, io.EOF) {
				session.Logger.Info("connection closed")
				return
			}
			session.Logger.Error(err.Error())
			return
		}
	}
}

func (s *Server) SendPacket(conn net.Conn, packet networking.Packet) error {
	packet.Write(conn)
	return nil
}

func (s *Server) ReadPacket(conn net.Conn) (networking.Packet, error) {
	readBuf := make([]byte, 1)
	n, err := conn.Read(readBuf)
	if err != nil || n != 1 {
		return nil, fmt.Errorf("could not read from connection: %w", err)
	}

	packetID := readBuf[0]

	switch packetID {
	case serverbound.C2SHandshakePacket:
		pkt := serverbound.Handshake{}
		pkt.Read(conn)
		return &pkt, nil
	case serverbound.C2SCreateMatch:
		pkt := serverbound.CreateMatch{}
		pkt.Read(conn)
		return &pkt, nil
	case serverbound.C2SJoinMatch:
		pkt := serverbound.JoinMatch{}
		pkt.Read(conn)
		return &pkt, nil
	default:
		return nil, ErrUnknownPacketType
	}
}

func (s *Server) addSession(session *Session) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions = append(s.sessions, session)
	return len(s.sessions) - 1
}

func (s *Server) createMatch(session *Session) (*game.Match, error) {
	id, err := s.generateMatchID()
	if err != nil {
		return nil, fmt.Errorf("generating match ID: %w", err)
	}

	match := game.NewMatch(id, session.conn)

	s.mu2.Lock()
	defer s.mu2.Unlock()
	s.matches[match.ID] = match

	return match, nil
}

func (s *Server) generateMatchID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("generating random UUID: %w", err)
	}
	return id.String(), nil
}
