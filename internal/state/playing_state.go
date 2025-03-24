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
	"github.com/typetrait/pingo/internal/math"
)

const (
	BallSize     = 8
	PaddleMargin = 55
)

type PlayingState struct {
	eventBus  event.EventBus
	Rules     game.Rules
	Bounds    game.Bounds
	PlayerOne *game.Player
	PlayerTwo *game.Player
	Ball      *game.Ball

	score map[*game.Player]int64
}

func NewPlayingState(eventBus event.EventBus, rules game.Rules, bounds game.Bounds, playerOne, playerTwo *game.Player, ball *game.Ball) *PlayingState {
	return &PlayingState{
		eventBus:  eventBus,
		Rules:     rules,
		Bounds:    bounds,
		PlayerOne: playerOne,
		PlayerTwo: playerTwo,
		Ball:      ball,
		score:     map[*game.Player]int64{},
	}
}

func (ps *PlayingState) Start() {
	ps.reset()

	ps.score[ps.PlayerOne] = 0
	ps.score[ps.PlayerTwo] = 0
}

func (ps *PlayingState) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf(
			"Player 1: %d pts. | Player 2: %d pts.",
			ps.score[ps.PlayerOne],
			ps.score[ps.PlayerTwo],
		),
	)

	// Player One
	vector.DrawFilledRect(
		screen,
		ps.PlayerOne.Paddle.Position.X,
		ps.PlayerOne.Paddle.Position.Y,
		ps.PlayerOne.Paddle.Size.X,
		ps.PlayerOne.Paddle.Size.Y,
		color.White,
		true,
	)

	// Player Two
	vector.DrawFilledRect(
		screen,
		ps.PlayerTwo.Paddle.Position.X,
		ps.PlayerTwo.Paddle.Position.Y,
		ps.PlayerTwo.Paddle.Size.X,
		ps.PlayerTwo.Paddle.Size.Y,
		color.White,
		true,
	)

	// Ball
	vector.DrawFilledRect(
		screen,
		ps.Ball.Position.X,
		ps.Ball.Position.Y,
		BallSize,
		BallSize,
		color.White,
		true,
	)
}

func (ps *PlayingState) Update(dt float32) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		ps.eventBus.Publish(&event.ExitGameEvent{})
		return
	}

	// Player One Input
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		ps.PlayerOne.Paddle.Position.Y += -7.2
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		ps.PlayerOne.Paddle.Position.Y += 7.2
	}

	// Player Two Input
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		ps.PlayerTwo.Paddle.Position.Y += -7.2
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		ps.PlayerTwo.Paddle.Position.Y += 7.2
	}

	ps.Ball.Position = math.Vector2fAdd(ps.Ball.Position, ps.Ball.Velocity)

	// Paddle One Collision
	if ps.Ball.Velocity.X < 0 {
		if ps.Ball.Position.X+BallSize >= ps.PlayerOne.Paddle.Position.X && ps.Ball.Position.X <= ps.PlayerOne.Paddle.Position.X+ps.PlayerOne.Paddle.Size.X {
			if ps.Ball.Position.Y+BallSize >= ps.PlayerOne.Paddle.Position.Y && ps.Ball.Position.Y <= ps.PlayerOne.Paddle.Position.Y+ps.PlayerOne.Paddle.Size.Y {
				ps.Ball.Velocity.X = -ps.Ball.Velocity.X
			}
		}
	}

	// Paddle Two Collision
	if ps.Ball.Velocity.X > 0 {
		if ps.Ball.Position.X+BallSize >= ps.PlayerTwo.Paddle.Position.X && ps.Ball.Position.X <= ps.PlayerTwo.Paddle.Position.X+ps.PlayerTwo.Paddle.Size.X {
			if ps.Ball.Position.Y+BallSize >= ps.PlayerTwo.Paddle.Position.Y && ps.Ball.Position.Y <= ps.PlayerTwo.Paddle.Position.Y+ps.PlayerTwo.Paddle.Size.Y {
				ps.Ball.Velocity.X = -ps.Ball.Velocity.X
			}
		}
	}

	// Bounds Collision
	if ps.Ball.Position.Y+BallSize >= float32(ps.Bounds.Height) || ps.Ball.Position.Y <= 0 {
		ps.Ball.Velocity.Y = -ps.Ball.Velocity.Y
	}

	if ps.Ball.Position.X+BallSize <= 0 {
		ps.onScore(ps.PlayerTwo)
		return
	} else if ps.Ball.Position.X >= float32(ps.Bounds.Width) {
		ps.onScore(ps.PlayerOne)
		return
	}
}

func (ps *PlayingState) onScore(player *game.Player) {
	ps.score[player]++

	if ps.score[player] >= ps.Rules.WinningScore {
		ps.score[ps.PlayerOne] = 0
		ps.score[ps.PlayerTwo] = 0
		ps.reset()
	}

	ps.reset()
}

func (ps *PlayingState) reset() {
	ps.Ball.Position = math.NewVector2f(float32(ps.Bounds.Width)/2, float32(ps.Bounds.Height)/2)
	ps.Ball.Velocity = math.NewVector2f(-1*2.25, 1*2.25)

	ps.PlayerOne.Paddle.Position = math.NewVector2f(
		PaddleMargin,
		(float32(ps.Bounds.Height)/2)-(ps.PlayerOne.Paddle.Size.Y/2),
	)

	ps.PlayerTwo.Paddle.Position = math.NewVector2f(
		float32(ps.Bounds.Width)-PaddleMargin-ps.PlayerTwo.Paddle.Size.X,
		(float32(ps.Bounds.Height)/2)-(ps.PlayerTwo.Paddle.Size.Y/2),
	)
}
