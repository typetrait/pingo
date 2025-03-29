package game

import "github.com/typetrait/pingo/internal/math"

const (
	BasePaddleRoughness = 0.25
)

type Paddle struct {
	Position  math.Vector2f
	Size      math.Vector2f
	Roughness float32
}

func NewPaddle(position, size math.Vector2f, roughness float32) *Paddle {
	return &Paddle{
		Position:  position,
		Size:      size,
		Roughness: roughness,
	}
}
