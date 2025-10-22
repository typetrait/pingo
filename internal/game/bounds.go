package game

type Bounds struct {
	Width  int32
	Height int32
}

func NewBounds(width, height int32) *Bounds {
	return &Bounds{
		Width:  width,
		Height: height,
	}
}
