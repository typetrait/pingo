package state

import (
	"context"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/typetrait/pingo/internal/event"
	"github.com/typetrait/pingo/internal/net"
)

type PlayerInWaiting struct {
	Name string
}

func NewPlayerInWaiting(name string) *PlayerInWaiting {
	return &PlayerInWaiting{}
}

type MatchmakingState struct {
	eventBus  event.EventBus
	authority net.Authority

	connected bool
}

func NewMatchmakingState(bus event.EventBus, authority net.Authority) *MatchmakingState {
	return &MatchmakingState{
		eventBus:  bus,
		authority: authority,
		connected: false,
	}
}

func (mms *MatchmakingState) Start() {
	ctx, _ := context.WithCancel(context.Background())
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
}
