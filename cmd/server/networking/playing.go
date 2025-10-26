package networking

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/typetrait/pingo/cmd/server/game"
	"github.com/typetrait/pingo/internal/math"
	"github.com/typetrait/pingo/internal/packet/clientbound"
	"github.com/typetrait/pingo/internal/packet/serverbound"
)

const (
	stringStatePlaying string = "STATE_PLAYING"
)

const (
	tickRate uint8 = 60 // 60 updates per second
)

type PlayingSessionState struct {
	session *Session
	Match   *game.Match
}

func (s *PlayingSessionState) Handle(ctx context.Context) error {
	play := &clientbound.Play{}
	err := s.session.server.SendPacket(s.session.conn, play)
	if err != nil {
		return fmt.Errorf("sending play packet: %w", err)
	}

	ticker := time.NewTicker(time.Second / time.Duration(tickRate))
	defer ticker.Stop()

	var wg sync.WaitGroup
	wg.Add(2)

	// Input from this session's connection
	go func() {
		for {
			select {
			case <-ctx.Done():
				wg.Done()
			default:
				p, err := s.session.server.ReadPacket(s.session.conn)
				if err != nil {
				}

				switch pkt := p.(type) {
				case *serverbound.PaddleMove:
					s.session.Logger.Debug("paddle move", "match_id", s.Match.ID, "y_pos", pkt.Y)
				}
			}
		}
	}()

	// TODO: Move this
	go func() {
		for {
			select {
			case <-ctx.Done():
				wg.Done()
			case <-ticker.C:
				s.session.Logger.Debug("game state", "match_id", s.Match.ID)
				gameState := &clientbound.GameState{
					PlayerOneScore: 0,
					PlayerTwoScore: 0,
					PlayerOnePos: math.Vector2f{
						X: 0,
						Y: 0,
					},
					PlayerTwoPos: math.Vector2f{
						X: 0,
						Y: 0,
					},
					BallPos: math.Vector2f{
						X: 0,
						Y: 0,
					},
				}
				err := s.session.server.SendPacket(s.session.conn, gameState)
				if err != nil {
					s.session.Logger.Error("error sending game state packet", err)
				}
			}
		}
	}()
	wg.Wait()

	// TODO: Need a different state for this
	nss := &NegotiateSessionState{
		session: s.session,
	}
	s.session.SetState(nss)

	return nil
}

func (s *PlayingSessionState) String() string {
	return stringStatePlaying
}
