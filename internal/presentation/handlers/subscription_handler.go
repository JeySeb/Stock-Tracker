package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"stock-tracker/internal/domain/entities"
	"stock-tracker/internal/domain/usecases"
	"stock-tracker/internal/infrastructure/middleware"
	"stock-tracker/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	subscriptionUC usecases.SubscriptionUseCase
	validator      *validator.Validate
	logger         logger.Logger
}

func NewSubscriptionHandler(subscriptionUC usecases.SubscriptionUseCase, logger logger.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionUC: subscriptionUC,
		validator:      validator.New(),
		logger:         logger,
	}
}

func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		h.logger.Error("Failed to get user ID from context", "error", err)
		h.respondWithError(w, r, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req usecases.PaymentSimulationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request body", "error", err)
		h.respondWithError(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	if err := h.validator.Struct(req); err != nil {
		h.logger.Error("Request validation failed", "error", err)
		h.respondWithError(w, r, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	// Validate subscription plan
	if !isValidPlan(req.Plan) {
		h.respondWithError(w, r, http.StatusBadRequest, "Invalid subscription plan")
		return
	}

	subscription, err := h.subscriptionUC.CreateSubscription(r.Context(), userID, req)
	if err != nil {
		h.logger.Error("Failed to create subscription", "error", err, "userID", userID)
		h.respondWithError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	h.respondWithJSON(w, r, http.StatusCreated, subscription)
}

func (h *SubscriptionHandler) SimulatePayment(w http.ResponseWriter, r *http.Request) {
	subscriptionID, err := h.getSubscriptionIDFromURL(r)
	if err != nil {
		h.logger.Error("Invalid subscription ID", "error", err)
		h.respondWithError(w, r, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	// Verify user has access to this subscription
	_, err = h.getUserIDFromContext(r)
	if err != nil {
		h.logger.Error("Failed to get user ID from context", "error", err)
		h.respondWithError(w, r, http.StatusUnauthorized, "Authentication required")
		return
	}

	if err := h.subscriptionUC.SimulatePayment(r.Context(), subscriptionID); err != nil {
		h.logger.Error("Payment simulation failed", "error", err, "subscriptionID", subscriptionID)
		h.respondWithError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	h.respondWithJSON(w, r, http.StatusOK, map[string]string{"message": "Payment processed successfully"})
}

// Helper methods
func (h *SubscriptionHandler) getUserIDFromContext(r *http.Request) (uuid.UUID, error) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user ID not found in context")
	}
	return userID, nil
}

func (h *SubscriptionHandler) getSubscriptionIDFromURL(r *http.Request) (uuid.UUID, error) {
	subscriptionIDStr := chi.URLParam(r, "id")
	return uuid.Parse(subscriptionIDStr)
}

func (h *SubscriptionHandler) respondWithError(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, map[string]string{"error": message})
}

func (h *SubscriptionHandler) respondWithJSON(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	render.Status(r, status)
	render.JSON(w, r, data)
}

func isValidPlan(plan entities.SubscriptionPlan) bool {
	return plan == entities.PlanMonthly || plan == entities.PlanYearly
}
