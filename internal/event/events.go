package event

const (
	EventStartGame Type = iota
)

type Type int32

type Event interface {
	Type() Type
}

type StartGameEvent struct {
}

func (e *StartGameEvent) Type() Type {
	return EventStartGame
}
