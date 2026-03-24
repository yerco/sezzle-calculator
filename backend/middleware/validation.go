package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"sezzle-calculator/calculator"
	"sezzle-calculator/models"
)

// contextKey is an unexported type to avoid context key collisions.
type contextKey struct{}

// unaryOps lists operations that only use operand a; b must be zero.
var unaryOps = map[string]bool{
	"sqrt":       true,
	"percentage": true,
}

// GetRequest retrieves the validated CalculationRequest stored by ValidateCalculation.
func GetRequest(ctx context.Context) (models.CalculationRequest, bool) {
	req, ok := ctx.Value(contextKey{}).(models.CalculationRequest)
	return req, ok
}

// ValidateCalculation decodes the JSON body, rejects unknown fields and trailing
// garbage, validates the operation name, and passes the request via context.
func ValidateCalculation(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		var req models.CalculationRequest
		if err := dec.Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.CalculationResponse{Error: "invalid request body: " + err.Error()})
			return
		}

		// Reject trailing garbage after the JSON object (e.g. `{...}JUNK`).
		if dec.More() {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.CalculationResponse{Error: "invalid request body: unexpected data after JSON object"})
			return
		}

		if _, ok := calculator.GetOperation(req.Operation); !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.CalculationResponse{
				Error: "unsupported operation: " + req.Operation,
			})
			return
		}

		// For unary operations, b must be zero — catch accidental misuse early.
		if unaryOps[req.Operation] && req.B != 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.CalculationResponse{
				Error: req.Operation + " is a unary operation; b must be 0",
			})
			return
		}

		ctx := context.WithValue(r.Context(), contextKey{}, req)
		next(w, r.WithContext(ctx))
	}
}
