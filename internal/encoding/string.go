package encoding

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// WriteVarString writes a uint32 length prefix (in the given byte order) followed by the raw bytes of s.
func WriteVarString(w io.Writer, order binary.ByteOrder, s string) error {
	n := len(s)
	if uint64(n) > math.MaxUint32 {
		return fmt.Errorf("string too long: %d bytes", n)
	}
	if err := binary.Write(w, order, uint32(n)); err != nil {
		return err
	}
	_, err := io.WriteString(w, s)
	return err
}

// ReadVarString reads a uint32 length prefix (in the given byte order) and then that many bytes.
func ReadVarString(r io.Reader, order binary.ByteOrder) (string, error) {
	var n uint32
	if err := binary.Read(r, order, &n); err != nil {
		return "", err
	}

	// Optional safety cap (e.g., 16 MiB) to avoid pathological allocations.
	const maxSize = 16 << 20
	if n > maxSize {
		return "", fmt.Errorf("declined to read string of size %d (> %d)", n, maxSize)
	}

	if n == 0 {
		return "", nil
	}

	buf := make([]byte, n)
	if _, err := io.ReadFull(r, buf); err != nil {
		return "", err
	}
	return string(buf), nil
}
