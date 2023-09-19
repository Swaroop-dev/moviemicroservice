package http

import (
	"encoding/json"
	"log"
	"net/http"

	"movieapp.com/movie/internal/controller/movie"

	"errors"
)

type Handler struct {
	ctrl *movie.Controller
}

func New(c *movie.Controller) *Handler {
	return &Handler{c}
}

func (h *Handler) GetMovieDetails(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	details, err := h.ctrl.Get(req.Context(), id)
	if err != nil && errors.Is(err, movie.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(details); err != nil {
		log.Printf("Response encode error: %v\n", err)
	}
}
