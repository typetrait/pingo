package serverbound

import (
	"encoding/binary"
	"io"

	"github.com/typetrait/pingo/internal/encoding"
)

type JoinMatch struct {
	id         uint8
	SessionID  uint64
	MatchID    string
	PlayerName string
}

func (p *JoinMatch) ID() uint8 {
	return C2SJoinMatch
}

func (p *JoinMatch) Read(reader io.Reader) {
	_ = binary.Read(reader, binary.LittleEndian, &p.SessionID)
	p.MatchID, _ = encoding.ReadVarString(reader, binary.LittleEndian)
	p.PlayerName, _ = encoding.ReadVarString(reader, binary.LittleEndian)
}

func (p *JoinMatch) Write(writer io.Writer) {
	_ = binary.Write(writer, binary.LittleEndian, p.ID())
	_ = binary.Write(writer, binary.LittleEndian, p.SessionID)
	_ = encoding.WriteVarString(writer, binary.LittleEndian, p.MatchID)
	_ = encoding.WriteVarString(writer, binary.LittleEndian, p.PlayerName)
}
