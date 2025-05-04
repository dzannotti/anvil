package uuid

import (
	"crypto/rand"
	"fmt"
	"io"
)

func New() string {
	uuid := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, uuid); err != nil {
		return ""
	}

	// Set version (4) and variant bits according to RFC 4122
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])
}
