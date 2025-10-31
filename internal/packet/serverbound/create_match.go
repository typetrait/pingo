package serverbound

import (
	"encoding/binary"
	"io"

	"github.com/typetrait/pingo/internal/encoding"
)

type CreateMatch struct {
	id         uint8
	PlayerName string
}

func (p *CreateMatch) ID() uint8 {
	return C2SCreateMatch
}

func (p *CreateMatch) Read(reader io.Reader) {
	p.PlayerName, _ = encoding.ReadVarString(reader, binary.LittleEndian)
}

func (p *CreateMatch) Write(writer io.Writer) {
	_ = binary.Write(writer, binary.LittleEndian, p.ID())
	_ = encoding.WriteVarString(writer, binary.LittleEndian, p.PlayerName)
}
