package clientbound

import (
	"encoding/binary"
	"io"

	"github.com/typetrait/pingo/internal/encoding"
)

type MatchCreated struct {
	id      uint8
	MatchID string
}

func (p *MatchCreated) ID() uint8 {
	return S2CMatchCreated
}

func (p *MatchCreated) Read(reader io.Reader) {
	// _ = binary.Read(reader, binary.LittleEndian, &p.id)
	p.MatchID, _ = encoding.ReadVarString(reader, binary.LittleEndian)
}

func (p *MatchCreated) Write(writer io.Writer) {
	_ = binary.Write(writer, binary.LittleEndian, p.ID())
	_ = encoding.WriteVarString(writer, binary.LittleEndian, p.MatchID)
}
