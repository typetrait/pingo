package networking

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

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
	Match   *Match
}

func (s *PlayingSessionState) Handle(ctx context.Context) error {
	play := &clientbound.Play{
		AdversaryName: s.Match.AdversaryPlayer(s.session).Name,
	}
	err := s.session.server.SendPacket(s.session.conn, play)
	if err != nil {
		return fmt.Errorf("sending play packet: %w", err)
	}

	ticker := time.NewTicker(time.Second / time.Duration(tickRate))
	defer ticker.Stop()

	var wg sync.WaitGroup
	wg.Add(2)

	s.Match.Start(ctx)

	// Input
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				p, err := s.session.server.ReadPacket(s.session.conn)
				if err != nil {
					s.session.Logger.Error("failed to read packet", err)
					if errors.Is(err, io.EOF) {
						// TODO: Cancel
						return
					}
					continue
				}

				switch pkt := p.(type) {
				case *serverbound.PaddleMove:
					s.session.Logger.Debug("paddle move", "match_id", s.Match.ID, "y_pos", pkt.Y)
					s.Match.Game.QueueInput(&MoveInput{
						PlayerID: s.session.ID,
						PosY:     pkt.Y,
					})
				}
			}
		}
	}()

	// Output
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case snap := <-s.Match.Game.Snapshots():
				gameState := &clientbound.GameState{
					PlayerOneScore: snap.P1Score,
					PlayerTwoScore: snap.P2Score,
					PlayerOnePosY:  snap.P1Y,
					PlayerTwoPosY:  snap.P2Y,
					BallPos:        snap.BallPos,
				}
				if err := s.session.server.SendPacket(s.session.conn, gameState); err != nil {
					s.session.Logger.Error("failed to send game state", err)
					// TODO: Cancel
					return
				}
			}
		}
	}()
	wg.Wait()

	return nil
}

func (s *PlayingSessionState) String() string {
	return stringStatePlaying
}
