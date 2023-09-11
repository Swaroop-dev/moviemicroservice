package httphandler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"movieapp.com/metadata/internal/controller/metadata"
	"movieapp.com/metadata/internal/repository"
)

type Handler struct {
	ctrl *metadata.Controller
}

func New(ctrl *metadata.Controller) *Handler {
	return &Handler{ctrl}
}

func (h *Handler) GetMetaData(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	res, err := h.ctrl.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrornotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Response encode error: %v\n", err)
	}
}
