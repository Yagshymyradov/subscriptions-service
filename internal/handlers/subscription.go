package handlers

import (
	"encoding/json"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/Yagshymyradov/subscriptions-service/internal/models"
	"github.com/Yagshymyradov/subscriptions-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Handler struct {
	svc    *service.SubscriptionService
	logger *zap.Logger
}

func New(svc *service.SubscriptionService, logger *zap.Logger) *Handler {
	return &Handler{svc: svc, logger: logger}
}

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", h.create)
		r.Get("/{id}", h.get)
		r.Get("/", h.list)
		r.Put("/{id}", h.update)
		r.Delete("/{id}", h.delete)
	})
	r.Get("/subscriptions/total", h.totalCost)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req models.Subscription
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := h.svc.Create(r.Context(), &req); err != nil {
		h.logger.Error("create", zap.Error(err))
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(req)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	sub, err := h.svc.Get(r.Context(), id)
	if err != nil {
		h.logger.Error("get", zap.Error(err))
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "no subscriptions found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sub); err != nil {
		h.logger.Error("get", zap.Error(err))
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}
	if _, err := uuid.Parse(userID); err != nil {
		http.Error(w, "invalid userID", http.StatusBadRequest)
		return
	}

	subs, err := h.svc.List(r.Context(), userID)
	if err != nil {
		h.logger.Error("list", zap.Error(err))
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(subs)

}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var req models.Subscription
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	req.ID = id
	if err := h.svc.Update(r.Context(), &req); err != nil {
		h.logger.Error("update", zap.Error(err))
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "no subscriptions found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		h.logger.Error("delete", zap.Error(err))
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "no subscriptions found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) totalCost(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}
	month := r.URL.Query().Get("month")
	monthInt, err := strconv.Atoi(month)
	if err != nil {
		http.Error(w, "invalid month", http.StatusBadRequest)
	}
	year := r.URL.Query().Get("year")
	if year == "" {
		http.Error(w, "year is required", http.StatusBadRequest)
		return
	}
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		http.Error(w, "invalid year", http.StatusBadRequest)
		return
	}
	if month == "" || year == "" {
		http.Error(w, "month and year are required", http.StatusBadRequest)
		return
	}
	serviceFilter := r.URL.Query().Get("serviceFilter")
	total, err := h.svc.TotalCost(r.Context(), userID, monthInt, yearInt, serviceFilter)
	if err != nil {
		h.logger.Error("totalCost", zap.Error(err))
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]int{"total": total}); err != nil {
		h.logger.Error("totalCost", zap.Error(err))
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
}
