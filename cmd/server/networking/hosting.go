package networking

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/typetrait/pingo/cmd/server/game"
	"github.com/typetrait/pingo/internal/packet/serverbound"
)

const (
	stringStateHostingMatch string = "STATE_HOSTING_MATCH"
)

type HostingMatchSessionState struct {
	session *Session
	Match   *game.Match
}

func (s *HostingMatchSessionState) Handle(ctx context.Context) error {
	p, err := s.session.server.ReadPacket(s.session.conn)
	if err != nil {
		return fmt.Errorf("reading packet: %w", err)
	}

	joinMatch, ok := p.(*serverbound.JoinMatch)
	if !ok {
		return fmt.Errorf("unexpected packet blurgh")
	}

	if joinMatch.MatchID != s.Match.ID {
		slog.Debug("got match id %q", joinMatch.MatchID)
		return fmt.Errorf("unexpected match ID")
	}

	secondPlayer := &game.Player{
		Name:  joinMatch.PlayerName,
		Score: 0,
	}

	slog.Info("player %q joined the match", secondPlayer.Name)

	pss := &PlayingSessionState{
		session: s.session,
		Match:   s.Match,
	}
	s.session.SetState(pss)

	return nil
}

func (s *HostingMatchSessionState) String() string {
	return stringStateHostingMatch
}
