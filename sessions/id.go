package sessions

import (
	"math/rand"

	"github.com/oklog/ulid"
)

var (
	entropy = rand.New(rand.NewSource(int64(ulid.Now())))
)

// NewID creates a new session id (ULID).
func NewID() string {
	id := ulid.MustNew(ulid.Now(), entropy)
	return id.String()
}
