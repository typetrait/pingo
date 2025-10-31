package networking

import (
	"context"
	"time"
)

const (
	stringStateHostingMatch string = "STATE_HOSTING_MATCH"
)

type HostingMatchSessionState struct {
	session *Session
	Match   *Match
}

func (s *HostingMatchSessionState) Handle(ctx context.Context) error {
	for !s.Match.Ready() {
		s.session.Logger.Info("waiting for player to join")
		time.Sleep(time.Second * 3)
	}

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
