package math

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

func Vector2fAdd(first, second Vector2f) Vector2f {
	return NewVector2f(
		first.X+second.X,
		first.Y+second.Y,
	)
}
