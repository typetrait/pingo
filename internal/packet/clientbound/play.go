package clientbound

import (
	"encoding/binary"
	"io"

	"github.com/typetrait/pingo/internal/encoding"
)

type Play struct {
	id            uint8
	AdversaryName string
}

func (p *Play) ID() uint8 {
	return S2CPlay
}

func (p *Play) Read(reader io.Reader) {
	p.AdversaryName, _ = encoding.ReadVarString(reader, binary.LittleEndian)
}

func (p *Play) Write(writer io.Writer) {
	_ = binary.Write(writer, binary.LittleEndian, p.ID())
	_ = encoding.WriteVarString(writer, binary.LittleEndian, p.AdversaryName)
}
