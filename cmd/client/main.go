package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/typetrait/pingo/internal/event"
	"github.com/typetrait/pingo/internal/state"
)

type Game struct {
	width  int32
	height int32
	state  state.GameState
}

func NewGame(width, height int32) *Game {
	return &Game{
		width:  width,
		height: height,
	}
}

func (g *Game) Run() {
	bus := event.NewEventBussin()

	bus.Register(event.EventSetGameState, func(ev event.Event) {
		changeStateEvent := ev.(*state.SetGameStateEvent)
		g.setState(changeStateEvent.State)
	})

	bus.Register(event.EventExitGame, func(ev event.Event) {
		ms := state.NewMenuState(bus)
		g.setState(ms)
	})

	bus.Register(event.EventStartGame, func(ev event.Event) {
		startGameEvent := ev.(*state.StartGameEvent)
		ps := state.NewPlayingState(
			startGameEvent.Logic,
			bus,
			startGameEvent.Rules,
			startGameEvent.Bounds,
		)
		g.setState(ps)
	})

	g.state = state.NewMenuState(bus)
	g.state.Start()

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Pingo")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	g.state.Update(0)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.state.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}

func (g *Game) setState(state state.GameState) {
	g.state = state
	g.state.Start()
}

func main() {
	g := NewGame(800, 600)
	g.Run()
}
