package networking

import (
	"context"
	"net"
)

type Match struct {
	ID        string
	Game      *Game
	host      net.Conn
	playerOne *Player
	playerTwo *Player
}

func NewMatch(id string, host net.Conn) *Match {
	return &Match{
		ID:   id,
		host: host,
	}
}

func (m *Match) PlayerOne() *Player {
	return m.playerOne
}

func (m *Match) PlayerTwo() *Player {
	return m.playerTwo
}

func (m *Match) Ball() *Ball {
	return m.Game.Ball
}

func (m *Match) SetPlayerOne(player *Player) {
	m.playerOne = player
}

func (m *Match) SetPlayerTwo(player *Player) {
	m.playerTwo = player
}

func (m *Match) Ready() bool {
	return m.playerOne != nil && m.playerTwo != nil
}

func (m *Match) Host() net.Conn {
	return m.host
}

func (m *Match) Start(ctx context.Context) {
	m.Game = NewGame(m.playerOne, m.playerTwo)
	go m.Game.Run(ctx, tickRate)
}

func (m *Match) SessionPlayer(session *Session) *Player {
	if session.ID == m.playerOne.ID {
		return m.playerOne
	} else if session.ID == m.playerTwo.ID {
		return m.playerTwo
	}
	return nil
}

func (m *Match) AdversaryPlayer(session *Session) *Player {
	if session.ID == m.playerOne.ID {
		return m.playerTwo
	} else if session.ID == m.playerTwo.ID {
		return m.playerOne
	}
	return nil
}
