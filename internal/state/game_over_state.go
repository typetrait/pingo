package state

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/typetrait/pingo/internal/event"
	"github.com/typetrait/pingo/internal/game"
)

type GameOverEvent struct {
	winner *game.Player
}

func NewGameOverEvent(winner *game.Player) *GameOverEvent {
	return &GameOverEvent{
		winner: winner,
	}
}

func (e *GameOverEvent) Type() event.Type {
	return event.EventGameOver
}

type GameOverState struct {
	eventBus     event.EventBus
	playingState *PlayingState
	winner       *game.Player
}

func NewGameOverState(eventBus event.EventBus, playingState *PlayingState, winner *game.Player) *GameOverState {
	return &GameOverState{
		eventBus:     eventBus,
		playingState: playingState,
		winner:       winner,
	}
}

func (gos *GameOverState) Start() {

}

func (gos *GameOverState) Draw(screen *ebiten.Image) {
	gos.playingState.Draw(screen)
	vector.DrawFilledRect(
		screen,
		0,
		0,
		800,
		600,
		color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 128,
		},
		true,
	)

	const (
		cw = 6
		ch = 16
	)

	txt := fmt.Sprintf("%s wins!", gos.winner.Name)
	textWidth := (cw * len(txt))
	textHeight := ch
	ebitenutil.DebugPrintAt(screen, txt, 400-(textWidth/2), 300-(textHeight/2))
}

func (gos *GameOverState) Update(dt float32) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		gos.eventBus.Publish(&event.StartGameEvent{})
		return
	}
}
