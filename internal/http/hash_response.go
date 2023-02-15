package http

import (
	"time"

	"hash-sv/internal/hash"
)

type HashRowResponse struct {
	Hash      hash.UUID `json:"hash"`
	CreatedAt time.Time `json:"created_at"`
}

func NewHashRowResponse(h hash.UUID, t time.Time) HashRowResponse {
	return HashRowResponse{Hash: h, CreatedAt: t}
}
