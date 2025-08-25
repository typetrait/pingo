package state

import (
	"fmt"
	"image/color"
	math2 "math"
	"math/rand/v2"
	"time"

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

	BallSpeed   = 8.5
	PaddleSpeed = 10

	ballKickoffDelay = 1 * time.Second
)

type SetGameStateEvent struct {
	State GameState
}

func NewSetGameStateEvent(state GameState) *SetGameStateEvent {
	return &SetGameStateEvent{
		State: state,
	}
}

func (e *SetGameStateEvent) Type() event.Type {
	return event.EventSetGameState
}

type PlayingState struct {
	eventBus  event.EventBus
	Rules     game.Rules
	Bounds    game.Bounds
	PlayerOne *game.Player
	PlayerTwo *game.Player
	Ball      *game.Ball

	score map[*game.Player]int64

	isRoundOver bool
	roundOverAt time.Time

	kickoffFrames int

	mode GameMode
}

func NewPlayingState(eventBus event.EventBus, rules game.Rules, bounds game.Bounds, playerOne, playerTwo *game.Player, ball *game.Ball) *PlayingState {
	return &PlayingState{
		eventBus:      eventBus,
		Rules:         rules,
		Bounds:        bounds,
		PlayerOne:     playerOne,
		PlayerTwo:     playerTwo,
		Ball:          ball,
		score:         map[*game.Player]int64{},
		isRoundOver:   false,
		kickoffFrames: int(float64(ballKickoffDelay.Seconds() * float64(ebiten.TPS()))),
	}
}

func (ps *PlayingState) Start() {
	ps.eventBus.Register(event.EventGameOver, func(ev event.Event) {
		gameOverEvent := ev.(*GameOverEvent)
		ps.eventBus.Publish(
			NewSetGameStateEvent(
				NewGameOverState(ps.eventBus, ps, gameOverEvent.winner),
			),
		)
	})

	ps.score[ps.PlayerOne] = 0
	ps.score[ps.PlayerTwo] = 0

	ps.reset()
	ps.kickoffFrames = int(float64(ballKickoffDelay.Seconds() * float64(ebiten.TPS())))
}

func (ps *PlayingState) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf(
			"%s: %d pts. | %s: %d pts.",
			ps.PlayerOne.Name,
			ps.score[ps.PlayerOne],
			ps.PlayerTwo.Name,
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
	// Instant Game Over - for testing
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ps.eventBus.Publish(
			NewGameOverEvent(
				ps.PlayerOne,
			),
		)
		return
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		ps.eventBus.Publish(&event.ExitGameEvent{})
		return
	}

	// Player One Input
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		ps.PlayerOne.Paddle.Position.Y += -PaddleSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		ps.PlayerOne.Paddle.Position.Y += PaddleSpeed
	}

	// Player Two Input
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		ps.PlayerTwo.Paddle.Position.Y += -PaddleSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		ps.PlayerTwo.Paddle.Position.Y += PaddleSpeed
	}

	// Ball movement logic
	if ps.kickoffFrames > 0 {
		ps.kickoffFrames--

		if ps.kickoffFrames <= 0 {
			ps.kickoffFrames = 0
			ps.ballKickoff()
		}
	} else {
		ps.Ball.Position = ps.Ball.Position.Add(
			ps.Ball.Velocity.MultiplyByScalar(BallSpeed),
		)
	}

	// Paddle One Ball Collision
	if ps.Ball.Velocity.X < 0 {
		if ps.Ball.Position.X+BallSize >= ps.PlayerOne.Paddle.Position.X && ps.Ball.Position.X <= ps.PlayerOne.Paddle.Position.X+ps.PlayerOne.Paddle.Size.X {
			if ps.Ball.Position.Y+BallSize >= ps.PlayerOne.Paddle.Position.Y && ps.Ball.Position.Y <= ps.PlayerOne.Paddle.Position.Y+ps.PlayerOne.Paddle.Size.Y {
				ps.Ball.Velocity.X = -ps.Ball.Velocity.X
			}
		}
	}

	// Paddle Two Ball Collision
	if ps.Ball.Velocity.X > 0 {
		if ps.Ball.Position.X+BallSize >= ps.PlayerTwo.Paddle.Position.X && ps.Ball.Position.X <= ps.PlayerTwo.Paddle.Position.X+ps.PlayerTwo.Paddle.Size.X {
			if ps.Ball.Position.Y+BallSize >= ps.PlayerTwo.Paddle.Position.Y && ps.Ball.Position.Y <= ps.PlayerTwo.Paddle.Position.Y+ps.PlayerTwo.Paddle.Size.Y {
				ps.Ball.Velocity.X = -ps.Ball.Velocity.X
			}
		}
	}

	// Paddle One Bounds Collision
	if ps.PlayerOne.Paddle.Position.Y <= 0 {
		ps.PlayerOne.Paddle.Position.Y = 0
	}

	if ps.PlayerOne.Paddle.Position.Y+ps.PlayerOne.Paddle.Size.Y >= float32(ps.Bounds.Height) {
		ps.PlayerOne.Paddle.Position.Y = float32(ps.Bounds.Height) - ps.PlayerTwo.Paddle.Size.Y
	}

	// Paddle Two Bounds Collision
	if ps.PlayerTwo.Paddle.Position.Y <= 0 {
		ps.PlayerTwo.Paddle.Position.Y = 0
	}

	if ps.PlayerTwo.Paddle.Position.Y+ps.PlayerTwo.Paddle.Size.Y >= float32(ps.Bounds.Height) {
		ps.PlayerTwo.Paddle.Position.Y = float32(ps.Bounds.Height) - ps.PlayerTwo.Paddle.Size.Y
	}

	// Ball Bounds Collision
	if ps.Ball.Position.Y+BallSize >= float32(ps.Bounds.Height) || ps.Ball.Position.Y <= 0 {
		ps.Ball.Velocity.Y = -ps.Ball.Velocity.Y
	}

	if !ps.isRoundOver {
		if ps.Ball.Position.X+BallSize <= 0 {
			ps.onScore(ps.PlayerTwo)
			return
		} else if ps.Ball.Position.X >= float32(ps.Bounds.Width) {
			ps.onScore(ps.PlayerOne)
			return
		}
	}
}

func (ps *PlayingState) onRoundStart() {
	ps.isRoundOver = false
	ps.resetBall()
}

func (ps *PlayingState) onScore(player *game.Player) {
	ps.score[player]++
	if ps.score[player] >= ps.Rules.WinningScore {
		ps.onGameOver(player)
		return
	}
	ps.kickoffFrames = int(float64(ballKickoffDelay.Seconds() * float64(ebiten.TPS())))
	ps.onRoundStart()
}

func (ps *PlayingState) onGameOver(winner *game.Player) {
	// for k, _ := range ps.score {
	// 	ps.score[k] = 0
	// }

	ps.eventBus.Publish(
		NewSetGameStateEvent(
			NewGameOverState(ps.eventBus, ps, winner),
		),
	)
}

func (ps *PlayingState) reset() {
	ps.resetBall()
	ps.resetPaddles()
}

func (ps *PlayingState) resetBall() {
	ps.Ball.Position = math.NewVector2f(float32(ps.Bounds.Width)/2, float32(ps.Bounds.Height)/2)
}

func (ps *PlayingState) resetPaddles() {
	ps.PlayerOne.Paddle.Position = math.NewVector2f(
		PaddleMargin,
		(float32(ps.Bounds.Height)/2)-(ps.PlayerOne.Paddle.Size.Y/2),
	)

	ps.PlayerTwo.Paddle.Position = math.NewVector2f(
		float32(ps.Bounds.Width)-PaddleMargin-ps.PlayerTwo.Paddle.Size.X,
		(float32(ps.Bounds.Height)/2)-(ps.PlayerTwo.Paddle.Size.Y/2),
	)
}

func (ps *PlayingState) ballKickoff() {
	ps.Ball.Velocity = ps.randomBallVelocity()
}

func (ps *PlayingState) randomBallVelocity() math.Vector2f {
	var angle float64
	if rand.IntN(2) == 0 {
		angle = rand.Float64()*math2.Pi/2 - math2.Pi/4
	} else {
		angle = rand.Float64()*math2.Pi/2 + 3*math2.Pi/4
	}

	x := math2.Cos(angle)
	y := math2.Sin(angle)

	velocity := math.NewVector2f(
		float32(x),
		float32(y),
	).Normalize()

	return velocity
}
