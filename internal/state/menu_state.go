package state

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/typetrait/pingo/assets"
	"github.com/typetrait/pingo/internal/event"
)

const (
	pressEnterLabelText = "Press ENTER to start playing..."
)

type MenuState struct {
	eventBus event.EventBus

	font      *text.GoTextFace
	logoImage *ebiten.Image
}

func NewMenuState(eventBus event.EventBus) *MenuState {
	return &MenuState{
		eventBus: eventBus,
	}
}

func (ms *MenuState) Start() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	ms.font = &text.GoTextFace{
		Source: s,
		Size:   16,
	}

	// ---

	logoImg, _, err := image.Decode(bytes.NewReader(assets.Pingo_png))
	if err != nil {
		log.Fatal(err)
	}

	ms.logoImage = ebiten.NewImageFromImage(logoImg)
}

func (ms *MenuState) Draw(screen *ebiten.Image) {
	halfLogoWidth := float64(ms.logoImage.Bounds().Size().X / 2)
	halfLogoHeight := float64(ms.logoImage.Bounds().Size().Y / 2)

	drawOptions := &text.DrawOptions{}
	w, h := text.Measure(pressEnterLabelText, ms.font, drawOptions.LineSpacing)

	logoDrawOptions := &ebiten.DrawImageOptions{}
	logoDrawOptions.GeoM.Translate(400-halfLogoWidth, 300-halfLogoHeight-h)
	screen.DrawImage(ms.logoImage, logoDrawOptions)

	drawOptions.GeoM.Translate(400-w/2, 300+h*2)
	text.Draw(screen, pressEnterLabelText, ms.font, drawOptions)
}

func (ms *MenuState) Update(dt float32) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		ms.eventBus.Publish(&event.StartGameEvent{})
		return
	}
}
