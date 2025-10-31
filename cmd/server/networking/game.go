package networking

import (
	"context"
	"time"

	"github.com/typetrait/pingo/internal/math"
)

type MoveInput struct {
	PlayerID uint64
	PosY     float64
}

type Snapshot struct {
	P1Y, P2Y         float64
	P1Score, P2Score uint64
	BallPos          math.Vector2f
}

type Player struct {
	ID    uint64
	Name  string
	PosY  float64
	Score uint64
}

type Ball struct {
	Position math.Vector2f
}

type Game struct {
	PlayerOne *Player
	PlayerTwo *Player
	Ball      *Ball
	inputs    chan *MoveInput
	snapshots chan *Snapshot
}

func NewGame(playerOne *Player, playerTwo *Player) *Game {
	g := &Game{
		PlayerOne: playerOne,
		PlayerTwo: playerTwo,
		Ball:      &Ball{},
		inputs:    make(chan *MoveInput),
		snapshots: make(chan *Snapshot),
	}

	return g
}

func (g *Game) Snapshots() <-chan *Snapshot {
	return g.snapshots
}

func (g *Game) QueueInput(input *MoveInput) {
	g.inputs <- input
}

func (g *Game) Run(ctx context.Context, tickRate uint8) {
	ticker := time.NewTicker(time.Second / time.Duration(tickRate))
	defer ticker.Stop()
	last := time.Now()
	vx, vy := 220.0, 180.0
	for {
		select {
		case <-ctx.Done():
			return
		case input := <-g.inputs:
			if input.PlayerID == g.PlayerOne.ID {
				_ = g.PlayerOne.PosY - input.PosY // delta
				g.PlayerOne.PosY = input.PosY
			} else if input.PlayerID == g.PlayerTwo.ID {
				_ = g.PlayerTwo.PosY - input.PosY // delta
				g.PlayerTwo.PosY = input.PosY
			}
		case <-ticker.C:
			dt := time.Since(last).Seconds()
			last = time.Now()
			g.Ball.Position.X += float32(vx * dt)
			g.Ball.Position.Y += float32(vy * dt)
			select {
			case g.snapshots <- &Snapshot{
				P1Y:     g.PlayerOne.PosY,
				P2Y:     g.PlayerTwo.PosY,
				P1Score: g.PlayerOne.Score,
				P2Score: g.PlayerTwo.Score,
				BallPos: g.Ball.Position,
			}:
			}
		default:
		}
	}
}
