package id

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func New() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	ts := time.Now().UnixNano()

	return fmt.Sprintf("%s-%x", hex.EncodeToString(b), ts)
}
