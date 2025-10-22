package serverbound

import (
	"encoding/binary"
	"io"

	"github.com/typetrait/pingo/internal/encoding"
)

type JoinMatch struct {
	id         uint8
	MatchID    string
	PlayerName string
}

func (p *JoinMatch) ID() uint8 {
	return p.id
}

func (p *JoinMatch) Read(reader io.Reader) {
	// _ = binary.Read(reader, binary.LittleEndian, &p.id)
	p.MatchID, _ = encoding.ReadVarString(reader, binary.LittleEndian)
	p.PlayerName, _ = encoding.ReadVarString(reader, binary.LittleEndian)
}

func (p *JoinMatch) Write(writer io.Writer) {
	_ = binary.Write(writer, binary.LittleEndian, p.id)
	_ = encoding.WriteVarString(writer, binary.LittleEndian, p.MatchID)
	_ = encoding.WriteVarString(writer, binary.LittleEndian, p.PlayerName)
}
