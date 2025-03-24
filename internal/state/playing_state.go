package state

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/typetrait/pingo/internal/game"
	"github.com/typetrait/pingo/internal/math"
)

const (
	BallSize = 8
)

type PlayingState struct {
	PlayerOne *game.Player
	PlayerTwo *game.Player
	Ball      *game.Ball

	score map[*game.Player]int32
}

func NewPlayingState(playerOne, playerTwo *game.Player, ball *game.Ball) *PlayingState {
	return &PlayingState{
		PlayerOne: playerOne,
		PlayerTwo: playerTwo,
		Ball:      ball,
		score:     map[*game.Player]int32{},
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

	// TODO: USE RESOLUTION/SCREEN SIZE
	// Bounds Collision
	if ps.Ball.Position.Y+BallSize >= 600 || ps.Ball.Position.Y <= 0 {
		ps.Ball.Velocity.Y = -ps.Ball.Velocity.Y
	}

	if ps.Ball.Position.X+BallSize <= 0 {
		ps.onScore(ps.PlayerTwo)
		return
	} else if ps.Ball.Position.X >= 800 {
		ps.onScore(ps.PlayerOne)
		return
	}
}

func (ps *PlayingState) onScore(player *game.Player) {
	ps.score[player]++
	ps.reset()
}

func (ps *PlayingState) reset() {
	ps.Ball.Position = math.NewVector2f(400, 300)
	ps.Ball.Velocity = math.NewVector2f(-1*2.25, 1*2.25)
}
