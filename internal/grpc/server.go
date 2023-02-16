package grpc

import (
	"context"

	"hash-sv/internal/hash"
)

type Server struct {
	UnimplementedHashServer

	service *hash.Service
}

func NewServer(service *hash.Service) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) Get(ctx context.Context, _ *Empty) (*HashRowResponse, error) {
	hashRow := s.service.Get(ctx)

	return &HashRowResponse{
		Hash:      hashRow.Hash.String(),
		CreatedAt: hashRow.CreatedAt.String(),
	}, nil
}
