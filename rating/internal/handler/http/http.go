package httpHandler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	rating "movieapp.com/rating/internal/controller"
	"movieapp.com/rating/pkg/model"
)

type Handler struct {
	ctrl *rating.Controller
}

func New(ctrl *rating.Controller) *Handler {
	return &Handler{ctrl}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	recordId := model.RecordId(r.FormValue("id"))

	if recordId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	recordType := model.RecordType(r.FormValue("type"))

	if recordType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		v, err := h.ctrl.GetAgrrementRating(r.Context(), recordId, recordType)

		if err != nil && errors.Is(err, rating.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(v); err != nil {
			log.Printf("error in encodeing %v ", v, err)
		}

	case http.MethodPut:
		userId := model.UserId(r.FormValue("userId"))
		value, err := strconv.ParseFloat(r.FormValue("value"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := h.ctrl.PutAgreementRating(r.Context(), recordId, recordType, &model.Rating{UserId: userId, Value: model.RatingValue(value)}); err != nil {
			log.Printf("Repository put error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

	default:
		w.WriteHeader(http.StatusBadGateway)
	}

}
