package math

import (
	math2 "math"
)

type Vector2f struct {
	X float32
	Y float32
}

var (
	ZeroVector2f = Vector2f{}
)

func NewVector2f(x, y float32) Vector2f {
	return Vector2f{
		X: x,
		Y: y,
	}
}

func (v Vector2f) Normalize() Vector2f {
	mag := float32(math2.Sqrt(float64(v.X)*float64(v.X) + float64(v.Y)*float64(v.Y)))
	if mag == 0 {
		return NewVector2f(0, 0)
	}
	return NewVector2f(v.X/mag, v.Y/mag)
}

func (v Vector2f) Add(other Vector2f) Vector2f {
	return Vector2fAdd(v, other)
}

func (v Vector2f) MultiplyByScalar(scalar float32) Vector2f {
	return Vector2fMultiplyByScalar(v, scalar)
}

func Vector2fAdd(first, second Vector2f) Vector2f {
	return NewVector2f(
		first.X+second.X,
		first.Y+second.Y,
	)
}

func Vector2fMultiplyByScalar(vec Vector2f, scalar float32) Vector2f {
	return NewVector2f(
		vec.X*scalar,
		vec.Y*scalar,
	)
}
