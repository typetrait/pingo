package event

const (
	EventStartGame Type = iota
	EventExitGame
	EventGameOver
	EventSetGameState
	EventMatchmaking
)

type Type int32

type Event interface {
	Type() Type
}

type ExitGameEvent struct {
}

func (e *ExitGameEvent) Type() Type {
	return EventExitGame
}
