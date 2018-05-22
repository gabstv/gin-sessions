package ledisdb

import (
	"encoding/binary"
	"fmt"
)

// K converts a string to a ledisdb key
func K(format string, a ...interface{}) []byte {
	return []byte(fmt.Sprintf(format, a...))
}

// ExportUint64 converts uint64 to bytes
func ExportUint64(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

// ExportInt64 converts int64 to bytes
func ExportInt64(v int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// ImportInt64 converts bytes to int64
func ImportInt64(v []byte) int64 {
	return int64(binary.BigEndian.Uint64(v))
}
