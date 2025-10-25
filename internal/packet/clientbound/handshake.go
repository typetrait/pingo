package clientbound

import (
	"encoding/binary"
	"io"
)

type Handshake struct {
	id uint8
}

func (p *Handshake) ID() uint8 {
	return S2CHandshake
}

func (p *Handshake) Read(reader io.Reader) {
	// _ = binary.Read(reader, binary.LittleEndian, &p.id)
}

func (p *Handshake) Write(writer io.Writer) {
	_ = binary.Write(writer, binary.LittleEndian, p.ID())
}
