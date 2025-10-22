package networking

import (
	"fmt"
	"sync"
	"time"

	"github.com/typetrait/pingo/cmd/server/game"
	"github.com/typetrait/pingo/internal/math"
	"github.com/typetrait/pingo/internal/packet/clientbound"
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

func (s *PlayingSessionState) Handle() error {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		//p, err := s.session.server.ReadPacket(s.session.conn)
		//if err != nil {
		//	return fmt.Errorf("reading packet: %w", err)
		//}
		//
		//joinMatch, ok := p.(*serverbound.JoinMatch)
		//if !ok {
		//	return fmt.Errorf("unexpected packet")
		//}
		//
		//if joinMatch.MatchID != s.Match.ID {
		//	return fmt.Errorf("unexpected match ID")
		//}
	}()

	timer := time.NewTimer(time.Duration(tickRate) * time.Second)
	go func() {
		defer wg.Done()

		select {
		case <-timer.C:
			fmt.Println("tick!")
			clientState := &clientbound.GameState{
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
			err := s.session.server.SendPacket(s.session.conn, clientState)
			if err != nil {
				fmt.Println("error sending game state packet:", err)
			}
		}
	}()

	wg.Wait()

	return nil
}

func (s *PlayingSessionState) String() string {
	return stringStatePlaying
}
