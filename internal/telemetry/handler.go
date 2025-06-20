package telemetry

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	service Service
	auth    func(http.Handler) http.Handler
}

func NewHandler(s Service, apiKeys ...string) *Handler {
	return &Handler{service: s, auth: AuthMiddleware(apiKeys)}
}

// WrapHandler applies the auth middleware to the given handler.
func (h *Handler) WrapHandler(handler http.Handler) http.Handler {
	return h.auth(handler)
}

func (h *Handler) GetScooter(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, r, http.StatusBadRequest, "invalid scooter id", err)
		return
	}

	scooter, err := h.service.GetScooter(r.Context(), id)
	if err != nil {
		h.Err(w, r, http.StatusNotFound, err.Error(), err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(scooter)
	if err != nil {
		h.Err(w, r, http.StatusInternalServerError, "response encoding error", err)
		return
	}
}

func (h *Handler) UpdateScooter(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, r, http.StatusBadRequest, "invalid scooter id", err)
		return
	}

	var s Scooter
	err = json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		h.Err(w, r, http.StatusBadRequest, "unmarshalable request body", err)
		return
	}

	s.ID = id

	err = h.service.UpdateScooter(r.Context(), s)
	if err != nil {
		h.Err(w, r, http.StatusInternalServerError, err.Error(), err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) FindScooters(w http.ResponseWriter, r *http.Request) {
	qry, err := NewQuery(r)
	if err != nil {
		h.Err(w, r, http.StatusBadRequest, "invalid area parameters", err)
		return
	}

	result, err := h.service.FindScooters(r.Context(), qry)
	if err != nil {
		h.Err(w, r, http.StatusInternalServerError, err.Error(), err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		h.Err(w, r, http.StatusInternalServerError, "response encoding error", err)
		return
	}
}

func (h *Handler) ReportEvent(w http.ResponseWriter, r *http.Request) {
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		h.Err(w, r, http.StatusBadRequest, "invalid event payload", err)
		return
	}

	err := h.service.ReportEvent(r.Context(), event)
	if err != nil {
		h.Err(w, r, http.StatusInternalServerError, err.Error(), err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Err(w http.ResponseWriter, r *http.Request, status int, msg string, err error) {
	http.Error(w, msg, status)
	clientID, _ := ClientID(r.Context())
	prefix := "handler error:"

	if clientID != "" {
		prefix = fmt.Sprintf("[%s] handler error:", clientID)
	}

	if err != nil {
		log.Printf("%s %s: %v", prefix, msg, err)
		return
	}

	log.Printf("%s %s", prefix, msg)
}
