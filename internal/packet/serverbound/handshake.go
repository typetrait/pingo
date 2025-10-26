package serverbound

import (
	"encoding/binary"
	"io"
)

type Handshake struct {
	id              uint8
	ProtocolVersion uint8
}

func (p *Handshake) ID() uint8 {
	return C2SHandshake
}

func (p *Handshake) Read(reader io.Reader) {
	// _ = binary.Read(reader, binary.LittleEndian, &p.id)
	_ = binary.Read(reader, binary.LittleEndian, &p.ProtocolVersion)
}

func (p *Handshake) Write(writer io.Writer) {
	_ = binary.Write(writer, binary.LittleEndian, p.ID())
	_ = binary.Write(writer, binary.LittleEndian, p.ProtocolVersion)
}
