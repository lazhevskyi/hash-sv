//go:generate mockgen -source=generator.go --destination=generator_mock.go -package=hash
package hash

import "github.com/google/uuid"

type UUID = uuid.UUID

type UUIDGenerator interface {
	Generate() UUID
}

type uuid4Generator struct{}

func NewUUID4Generator() UUIDGenerator {
	return &uuid4Generator{}
}

func (g uuid4Generator) Generate() UUID {
	return uuid.New()
}
