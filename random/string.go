package random

import (
	"crypto/rand"
	"fmt"
)

// String -
func String(size int) string {
	b := make([]byte, size/2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
