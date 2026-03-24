package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"sezzle-calculator/calculator"
	"sezzle-calculator/middleware"
	"sezzle-calculator/models"
)

// Handler holds the shared dependencies for all calculator endpoints.
type Handler struct {
	history *calculator.History
}

func NewHandler(history *calculator.History) *Handler {
	return &Handler{history: history}
}

// Calculate handles POST /calculate.
// It relies on the ValidateCalculation middleware having already decoded and
// validated the request — acting as the Originator that creates a HistoryEntry.
func (h *Handler) Calculate(w http.ResponseWriter, r *http.Request) {
	req, ok := middleware.GetRequest(r.Context())
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.CalculationResponse{Error: "internal error: missing request context"})
		return
	}

	op, _ := calculator.GetOperation(req.Operation)
	result, err := calculator.NewCalculator(op).Compute(req.A, req.B)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(models.CalculationResponse{Error: err.Error()})
		return
	}

	// Save a memento of this successful calculation.
	h.history.Save(calculator.HistoryEntry{
		A:         req.A,
		B:         req.B,
		Operation: req.Operation,
		Result:    result,
		Timestamp: time.Now().UTC(),
	})

	json.NewEncoder(w).Encode(models.CalculationResponse{Result: &result})
}

// History handles GET /history — returns all stored calculation mementos.
func (h *Handler) History(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	entries := h.history.Entries()
	// Return an empty array rather than null when there are no entries.
	if entries == nil {
		entries = []calculator.HistoryEntry{}
	}
	json.NewEncoder(w).Encode(entries)
}
