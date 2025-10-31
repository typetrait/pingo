package clientbound

import (
	"encoding/binary"
	"io"

	"github.com/typetrait/pingo/internal/math"
)

type GameState struct {
	id             uint8
	PlayerOneScore uint64
	PlayerTwoScore uint64
	PlayerOnePosY  float64
	PlayerTwoPosY  float64
	BallPos        math.Vector2f
}

func (p *GameState) ID() uint8 {
	return S2CGameState
}

func (p *GameState) Read(reader io.Reader) {
	_ = binary.Read(reader, binary.LittleEndian, &p.PlayerOneScore)
	_ = binary.Read(reader, binary.LittleEndian, &p.PlayerTwoScore)
	_ = binary.Read(reader, binary.LittleEndian, &p.PlayerOnePosY)
	_ = binary.Read(reader, binary.LittleEndian, &p.PlayerTwoPosY)
	_ = binary.Read(reader, binary.LittleEndian, &p.BallPos.X)
	_ = binary.Read(reader, binary.LittleEndian, &p.BallPos.Y)
}

func (p *GameState) Write(writer io.Writer) {
	_ = binary.Write(writer, binary.LittleEndian, p.ID())
	_ = binary.Write(writer, binary.LittleEndian, p.PlayerOneScore)
	_ = binary.Write(writer, binary.LittleEndian, p.PlayerTwoScore)
	_ = binary.Write(writer, binary.LittleEndian, p.PlayerOnePosY)
	_ = binary.Write(writer, binary.LittleEndian, p.PlayerTwoPosY)
	_ = binary.Write(writer, binary.LittleEndian, p.BallPos.X)
	_ = binary.Write(writer, binary.LittleEndian, p.BallPos.Y)
}
