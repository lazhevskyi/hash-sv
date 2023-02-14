package hash

import "github.com/google/uuid"

type UUID = uuid.UUID

type UUIDGenerator interface {
	Generate() UUID
}

type generator struct{}

func NewGenerator() UUIDGenerator {
	return &generator{}
}

func (g generator) Generate() UUID {
	return uuid.New()
}
