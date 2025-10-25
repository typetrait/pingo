package serverbound

import (
	"encoding/binary"
	"io"
)

type PaddleMove struct {
	id uint8
	Y  float32
}

func (p *PaddleMove) ID() uint8 {
	return C2SPaddleMovePacket
}

func (p *PaddleMove) Read(reader io.Reader) {
	// _ = binary.Read(reader, binary.LittleEndian, &p.id)
	_ = binary.Read(reader, binary.LittleEndian, &p.Y)
}

func (p *PaddleMove) Write(writer io.Writer) {
	_ = binary.Write(writer, binary.LittleEndian, p.ID())
	_ = binary.Write(writer, binary.LittleEndian, p.Y)
}
