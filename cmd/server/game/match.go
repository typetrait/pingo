package game

import "net"

type Match struct {
	ID     string
	Ready  bool
	host   net.Conn
	player net.Conn
}

func NewMatch(id string, host net.Conn) *Match {
	return &Match{
		ID:     id,
		Ready:  false,
		host:   host,
		player: nil,
	}
}

func (m *Match) SetPlayer(player net.Conn) {
	m.player = player
	m.Ready = true
}
