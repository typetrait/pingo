package networking

import (
	"context"
	"fmt"

	"github.com/typetrait/pingo/cmd/server/game"
	"github.com/typetrait/pingo/internal/packet/clientbound"
	"github.com/typetrait/pingo/internal/packet/serverbound"
)

const (
	stringStateNegotiate string = "STATE_NEGOTIATE"
)

type NegotiateSessionState struct {
	session *Session
}

func (s *NegotiateSessionState) Handle(ctx context.Context) error {
	// Receive and handle Handshake
	p, err := s.session.server.ReadPacket(s.session.conn)
	if err != nil {
		return fmt.Errorf("reading packet: %w", err)
	}

	handshake, ok := p.(*serverbound.Handshake)
	if !ok {
		return fmt.Errorf("unexpected packet")
	}

	err = s.handleHandshake(handshake)
	if err != nil {
		return fmt.Errorf("handshake failed: %w", err)
	}

	// Receive and handle subsequent allowed packets
	p, err = s.session.server.ReadPacket(s.session.conn)
	if err != nil {
		return fmt.Errorf("reading packet: %w", err)
	}

	switch pkt := p.(type) {
	case *serverbound.CreateMatch:
		m, err := s.handleCreateMatch(pkt)
		if err != nil {
			return fmt.Errorf("creating match: %w", err)
		}
		hms := &HostingMatchSessionState{
			session: s.session,
			Match:   m,
		}
		s.session.SetState(hms)
	case *serverbound.JoinMatch:
		err := s.handleJoinMatch(pkt)
		if err != nil {
			return fmt.Errorf("joining match: %w", err)
		}
		// TODO: Set State
	default:
		return fmt.Errorf("unexpected packet")
	}

	return nil
}

func (s *NegotiateSessionState) handleHandshake(p *serverbound.Handshake) error {
	s.session.Logger.Debug("handshake request received")
	if p.ProtocolVersion != protocolVersion {
		return fmt.Errorf("protocol version mismatch")
	}

	hs := &clientbound.Handshake{}
	err := s.session.server.SendPacket(s.session.conn, hs)
	if err != nil {
		return fmt.Errorf("sending client bound handshake packet: %w", err)
	}

	s.session.Logger.Info("handshake complete")
	return nil
}

func (s *NegotiateSessionState) handleCreateMatch(p *serverbound.CreateMatch) (*game.Match, error) {
	s.session.Logger.Debug("match creation request received")

	s.session.Logger.Debug("creating match")
	match, err := s.session.server.createMatch()
	if err != nil {
		return nil, fmt.Errorf("creating match: %w", err)
	}

	s.session.Logger.Info("match created", "id", match.ID)

	mc := &clientbound.MatchCreated{
		MatchID: match.ID,
	}

	err = s.session.server.SendPacket(s.session.conn, mc)
	if err != nil {
		return nil, fmt.Errorf("sending match created packet: %w", err)
	}

	return match, nil
}

func (s *NegotiateSessionState) handleJoinMatch(p *serverbound.JoinMatch) error {
	s.session.Logger.Debug("match join request received", "player", p.PlayerName)
	_, ok := s.session.server.matches[p.MatchID]
	if !ok {
		return fmt.Errorf("match not found")
	}
	return nil
}

func (s *NegotiateSessionState) String() string {
	return stringStateNegotiate
}
