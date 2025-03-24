package event

const (
	EventStartGame Type = iota
	EventExitGame
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

type ExitGameEvent struct {
}

func (e *ExitGameEvent) Type() Type {
	return EventExitGame
}
