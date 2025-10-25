package clientbound

import (
	"encoding/binary"
	"io"
)

type Play struct {
	id uint8
}

func (p *Play) ID() uint8 {
	return S2CPlay
}

func (p *Play) Read(reader io.Reader) {
}

func (p *Play) Write(writer io.Writer) {
	_ = binary.Write(writer, binary.LittleEndian, p.ID())
}
