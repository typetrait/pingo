package game

import "github.com/typetrait/pingo/internal/math"

type Ball struct {
	Position math.Vector2f
	Velocity math.Vector2f
}

func NewBall(position, velocity math.Vector2f) *Ball {
	return &Ball{
		Position: position,
		Velocity: velocity,
	}
}
