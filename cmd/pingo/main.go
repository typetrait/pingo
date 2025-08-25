package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/typetrait/pingo/internal/event"
	"github.com/typetrait/pingo/internal/game"
	"github.com/typetrait/pingo/internal/math"
	"github.com/typetrait/pingo/internal/net"
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

	bus.Register(event.EventSetGameState, func(event event.Event) {
		changeStateEvent := event.(*state.SetGameStateEvent)
		g.state = changeStateEvent.State
		g.state.Start()
	})

	bus.Register(event.EventExitGame, func(event event.Event) {
		g.state = state.NewMenuState(bus)
		g.state.Start()
	})

	bus.Register(event.EventStartGame, func(ev event.Event) {
		startGameEvent := ev.(*state.StartGameEvent)
		switch startGameEvent.Mode {
		case state.GameModeSinglePlayer:
			paddleWidth := float32(g.width) * 0.01
			paddleHeight := float32(g.height) * 0.25
			paddleMargin := float32(25)

			ballWidth := 8
			ballHeight := 8

			ball := game.NewBall(
				math.NewVector2f(
					float32(g.width)/2-(float32(ballWidth)/2), float32(g.height)/2-(float32(ballHeight)/2),
				),
				math.ZeroVector2f,
			)

			playerOne := game.NewPlayer(
				"Player One",
				game.NewPaddle(
					math.NewVector2f(paddleMargin, float32(g.height)/2-paddleHeight/2),
					math.NewVector2f(paddleWidth, paddleHeight),
					game.BasePaddleRoughness,
				),
			)

			playerTwo := game.NewPlayer(
				"Player Two",
				game.NewPaddle(
					math.NewVector2f((float32(g.width)-paddleWidth)-paddleMargin, float32(g.height)/2-paddleHeight/2),
					math.NewVector2f(paddleWidth, paddleHeight),
					game.BasePaddleRoughness,
				),
			)

			rules := game.NewRules(5)

			w, h := g.Layout(800, 600)
			bounds := game.NewBounds(
				int32(w),
				int32(h),
			)

			ps := state.NewPlayingState(
				bus,
				rules,
				bounds,
				playerOne,
				playerTwo,
				ball,
			)

			g.state = ps
			g.state.Start()

		case state.GameModeMultiPlayer:
			authority := net.NewServerAuthority()
			mms := state.NewMatchmakingState(bus, authority)

			g.state = mms
			g.state.Start()
		}
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

func main() {
	g := NewGame(800, 600)
	g.Run()
}
