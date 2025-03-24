package game

import "github.com/typetrait/pingo/internal/math"

type Paddle struct {
	Position math.Vector2f
	Size     math.Vector2f
}

func NewPaddle(Position, Size math.Vector2f) *Paddle {
	return &Paddle{
		Position: Position,
		Size:     Size,
	}
}
