package serverbound

import (
	"encoding/binary"
	"io"
)

type CreateMatch struct {
	id uint8
}

func (p *CreateMatch) ID() uint8 {
	return C2SCreateMatch
}

func (p *CreateMatch) Read(reader io.Reader) {
	// _ = binary.Read(reader, binary.LittleEndian, &p.id)
}

func (p *CreateMatch) Write(writer io.Writer) {
	_ = binary.Write(writer, binary.LittleEndian, p.ID())
}
