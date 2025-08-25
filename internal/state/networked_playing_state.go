package state

import "github.com/hajimehoshi/ebiten/v2"

type NetworkedPlayingState struct {
}

func NewNetworkedPlayingState() *NetworkedPlayingState {
	return &NetworkedPlayingState{}
}

func (nps *NetworkedPlayingState) Start() {

}

func (nps *NetworkedPlayingState) Draw(screen *ebiten.Image) {

}

func (nps *NetworkedPlayingState) Update(dt float32) {

}
