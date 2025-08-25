package state

import (
	"github.com/typetrait/pingo/internal/event"
)

type GameMode int

const (
	GameModeSinglePlayer GameMode = iota
	GameModeMultiPlayer
)

type StartGameEvent struct {
	Mode GameMode
}

func NewStartGameEvent(mode GameMode) StartGameEvent {
	return StartGameEvent{
		Mode: mode,
	}
}

func (e *StartGameEvent) Type() event.Type {
	return event.EventStartGame
}
