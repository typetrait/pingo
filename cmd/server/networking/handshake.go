package networking

import (
	"fmt"
	"log"

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

func (s *NegotiateSessionState) Handle() error {
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

	switch p.(type) {
	case *serverbound.CreateMatch:
		cm := p.(*serverbound.CreateMatch)
		m, err := s.handleCreateMatch(cm)
		if err != nil {
			return fmt.Errorf("creating match: %w", err)
		}
		hms := &HostingMatchSessionState{
			session: s.session,
			Match:   m,
		}
		s.session.SetState(hms)
	}

	return nil
}

func (s *NegotiateSessionState) handleHandshake(p *serverbound.Handshake) error {
	log.Println("got handshake packet")
	if p.ProtocolVersion != protocolVersion {
		return fmt.Errorf("protocol version mismatch")
	}

	hs := &clientbound.Handshake{}
	err := s.session.server.SendPacket(s.session.conn, hs)
	if err != nil {
		return fmt.Errorf("sending client bound handshake packet: %w", err)
	}

	log.Println("handshake complete")
	return nil
}

func (s *NegotiateSessionState) handleCreateMatch(p *serverbound.CreateMatch) (*game.Match, error) {
	log.Println("match creation requested")

	log.Println("creating match")
	match, err := s.session.server.createMatch()
	if err != nil {
		return nil, fmt.Errorf("creating match: %w", err)
	}

	log.Printf("match with id %s created", match.ID)

	mc := &clientbound.MatchCreated{
		MatchID: match.ID,
	}

	log.Println("notifying client of newly created match")

	err = s.session.server.SendPacket(s.session.conn, mc)
	if err != nil {
		return nil, fmt.Errorf("sending match created packet: %w", err)
	}

	return match, nil
}

func (s *NegotiateSessionState) String() string {
	return stringStateNegotiate
}
