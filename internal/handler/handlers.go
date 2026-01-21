package handler

import (
	"context"
	"encoding/json"
	"github.com/bdzhalalov/resilient-scatter-gather/internal/service"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Handler struct {
	service *service.MockServices
	logger  *logrus.Logger
}

type serviceResponse[T any] struct {
	value T
	err   error
}

func New(logger *logrus.Logger, service *service.MockServices) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	userCh := make(chan serviceResponse[string], 1)
	accessCh := make(chan serviceResponse[bool], 1)
	memoryCh := make(chan serviceResponse[string], 1)

	go h.getUserInfo(ctx, userCh)

	go h.checkAccess(ctx, accessCh)

	go h.getMemoryContext(ctx, memoryCh)

	var (
		user      string
		access    bool
		memoryCtx *string
	)

	for i := 0; i < 2; i++ {
		select {
		case res := <-userCh:
			if res.err != nil {
				http.Error(w, res.err.Error(), http.StatusInternalServerError)
				return
			}
			user = res.value

		case res := <-accessCh:
			if res.err != nil || !res.value {
				http.Error(w, res.err.Error(), http.StatusInternalServerError)
				return
			}
			access = res.value

		case <-ctx.Done():
			http.Error(w, "timeout", http.StatusInternalServerError)
			return
		}
	}

	select {
	case res := <-memoryCh:
		if res.err == nil {
			memoryCtx = &res.value
		}
	case <-ctx.Done():
		h.logger.Debug("response without non-critical service")
	}

	response := map[string]any{
		"user":   user,
		"access": access,
	}

	if memoryCtx != nil {
		response["memory"] = *memoryCtx
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) getUserInfo(ctx context.Context, userCh chan<- serviceResponse[string]) {
	u, err := h.service.GetUser(ctx)
	userCh <- serviceResponse[string]{u, err}
}

func (h *Handler) checkAccess(ctx context.Context, accessCh chan<- serviceResponse[bool]) {
	ok, err := h.service.CheckAccess(ctx)
	accessCh <- serviceResponse[bool]{ok, err}
}

func (h *Handler) getMemoryContext(ctx context.Context, memoryCh chan<- serviceResponse[string]) {
	v, err := h.service.GetContext(ctx)
	memoryCh <- serviceResponse[string]{v, err}
}
