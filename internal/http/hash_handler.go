package http

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"hash-sv/internal/hash"
)

type hashHandler struct {
	service *hash.Service
	logger  *zap.Logger
}

func NewHashHandler(
	service *hash.Service,
	logger *zap.Logger,
) http.Handler {
	return &hashHandler{
		service: service,
		logger:  logger,
	}
}

func (h *hashHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	row := h.service.Get()

	err := json.NewEncoder(w).Encode(row)
	if err != nil {
		h.logger.Error("failed to write response:", zap.Error(err))
	}
}
