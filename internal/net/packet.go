package net

import (
	"encoding/binary"
	"io"
)

type Packet interface {
	Read(reader io.Reader)
	Write(writer io.Writer)
}

type PaddleMovePacket struct {
	Y float32
}

func (p *PaddleMovePacket) Read(reader io.Reader) {
	_ = binary.Read(reader, binary.LittleEndian, &p.Y)
}

func (p *PaddleMovePacket) Write(writer io.Writer) {
	_ = binary.Write(writer, binary.LittleEndian, p.Y)
}
