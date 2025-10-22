package networking

import (
	"io"
)

type Packet interface {
	ID() uint8
	Read(reader io.Reader)
	Write(writer io.Writer)
}
