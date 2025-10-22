package state

import (
	"context"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/typetrait/pingo/internal/event"
	"github.com/typetrait/pingo/internal/game"
	"github.com/typetrait/pingo/internal/networking"
)

type PlayerInWaiting struct {
	Name string
}

func NewPlayerInWaiting(name string) *PlayerInWaiting {
	return &PlayerInWaiting{}
}

type MatchmakingState struct {
	eventBus  event.EventBus
	authority networking.Authority

	connected bool
}

func NewMatchmakingState(bus event.EventBus, authority networking.Authority) *MatchmakingState {
	return &MatchmakingState{
		eventBus:  bus,
		authority: authority,
		connected: false,
	}
}

func (mms *MatchmakingState) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := mms.authority.Connect(ctx); err != nil {
		os.Exit(1)
	}
	mms.connected = true
}

func (mms *MatchmakingState) Draw(screen *ebiten.Image) {
	if mms.connected {
		ebitenutil.DebugPrint(screen, "connected")
	} else {
		ebitenutil.DebugPrint(screen, "connecting...")
	}
}

func (mms *MatchmakingState) Update(dt float32) {
	logic := &NetGameLogic{}
	rules := game.NewRules(5)
	bounds := game.NewBounds(
		int32(800),
		int32(600),
	)

	if mms.connected {
		time.Sleep(3 * time.Second)
		mms.eventBus.Publish(
			NewStartGameEvent(
				logic,
				GameModeMultiPlayer,
				rules,
				bounds,
			),
		)
	}
}

func (mms *MatchmakingState) Type() event.Type {
	return event.EventMatchmaking
}
