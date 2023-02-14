package hash

import (
	"time"
)

type Row struct {
	Hash      UUID      `json:"hash"`
	CreatedAt time.Time `json:"created_at"`
}

func NewRow(hash UUID) Row {
	return Row{
		Hash:      hash,
		CreatedAt: time.Now(),
	}
}
