package clientbound

import (
	"encoding/binary"
	"io"

	"github.com/typetrait/pingo/internal/math"
)

type GameState struct {
	id           uint8
	PlayerOnePos math.Vector2f
	PlayerTwoPos math.Vector2f
	BallPos      math.Vector2f
}

func (p *GameState) ID() uint8 {
	return p.id
}

func (p *GameState) Read(reader io.Reader) {
	_ = binary.Read(reader, binary.LittleEndian, &p.id)
}

func (p *GameState) Write(writer io.Writer) {
	_ = binary.Write(writer, binary.LittleEndian, p.id)
	_ = binary.Write(writer, binary.LittleEndian, p.PlayerOnePos.X)
	_ = binary.Write(writer, binary.LittleEndian, p.PlayerOnePos.Y)
	_ = binary.Write(writer, binary.LittleEndian, p.PlayerTwoPos.X)
	_ = binary.Write(writer, binary.LittleEndian, p.PlayerTwoPos.Y)
	_ = binary.Write(writer, binary.LittleEndian, p.BallPos.X)
	_ = binary.Write(writer, binary.LittleEndian, p.BallPos.Y)
}
