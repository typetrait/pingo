package state

import "github.com/hajimehoshi/ebiten/v2"

type GameState interface {
	Start()
	Draw(screen *ebiten.Image)
	Update(dt float32)
}
