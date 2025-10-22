package state

import (
	"github.com/typetrait/pingo/internal/event"
	"github.com/typetrait/pingo/internal/game"
)

type GameMode int

const (
	GameModeSinglePlayer GameMode = iota
	GameModeMultiPlayer
)

type StartGameEvent struct {
	Logic  GameLogic
	Mode   GameMode
	Rules  *game.Rules
	Bounds *game.Bounds
}

func NewStartGameEvent(logic GameLogic, mode GameMode, rules *game.Rules, bounds *game.Bounds) *StartGameEvent {
	return &StartGameEvent{
		Logic:  logic,
		Mode:   mode,
		Rules:  rules,
		Bounds: bounds,
	}
}

func (e *StartGameEvent) Type() event.Type {
	return event.EventStartGame
}

// ---

type MatchmakingEvent struct {
}

func NewMatchmakingEvent() *MatchmakingEvent {
	return &MatchmakingEvent{}
}

func (e *MatchmakingEvent) Type() event.Type {
	return event.EventMatchmaking
}
