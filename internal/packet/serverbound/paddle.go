package serverbound

import (
	"encoding/binary"
	"io"
)

type PaddleMove struct {
	ID uint8
	Y  float32
}

func NewPaddleMove(y float32) *PaddleMove {
	return &PaddleMove{
		ID: C2SPaddleMovePacket,
		Y:  y,
	}
}

func (p *PaddleMove) Read(reader io.Reader) {
	_ = binary.Read(reader, binary.LittleEndian, &p.ID)
	_ = binary.Read(reader, binary.LittleEndian, &p.Y)
}

func (p *PaddleMove) Write(writer io.Writer) {
	_ = binary.Write(writer, binary.LittleEndian, p.ID)
	_ = binary.Write(writer, binary.LittleEndian, p.Y)
}
