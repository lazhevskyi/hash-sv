package hash

import (
	"time"
)

type Row struct {
	Hash      UUID
	CreatedAt time.Time
}

func NewRow(hash UUID) Row {
	return Row{
		Hash:      hash,
		CreatedAt: time.Now(),
	}
}
