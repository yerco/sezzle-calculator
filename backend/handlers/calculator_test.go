package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"sezzle-calculator/calculator"
	"sezzle-calculator/handlers"
	"sezzle-calculator/middleware"
	"sezzle-calculator/models"
)

// post sends a POST /calculate request through the middleware+handler chain.
func post(t *testing.T, h *handlers.Handler, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/calculate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	middleware.ValidateCalculation(h.Calculate).ServeHTTP(rr, req)
	return rr
}

func decodeResponse(t *testing.T, rr *httptest.ResponseRecorder) models.CalculationResponse {
	t.Helper()
	var resp models.CalculationResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return resp
}

func TestCalculateHandler_Addition(t *testing.T) {
	h := handlers.NewHandler(calculator.NewHistory())
	rr := post(t, h, `{"a":3,"b":4,"operation":"add"}`)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	resp := decodeResponse(t, rr)
	if resp.Error != "" {
		t.Fatalf("unexpected error: %s", resp.Error)
	}
	if resp.Result == nil || *resp.Result != 7 {
		t.Errorf("expected result 7, got %v", resp.Result)
	}
}

func TestCalculateHandler_ZeroResult(t *testing.T) {
	h := handlers.NewHandler(calculator.NewHistory())
	rr := post(t, h, `{"a":0,"b":0,"operation":"add"}`)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	resp := decodeResponse(t, rr)
	if resp.Result == nil || *resp.Result != 0 {
		t.Errorf("expected result 0, got %v", resp.Result)
	}
}

func TestCalculateHandler_DivisionByZero(t *testing.T) {
	h := handlers.NewHandler(calculator.NewHistory())
	rr := post(t, h, `{"a":5,"b":0,"operation":"divide"}`)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rr.Code)
	}
	resp := decodeResponse(t, rr)
	if resp.Error == "" {
		t.Error("expected error message for division by zero")
	}
}

func TestCalculateHandler_SqrtNegative(t *testing.T) {
	h := handlers.NewHandler(calculator.NewHistory())
	rr := post(t, h, `{"a":-4,"b":0,"operation":"sqrt"}`)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rr.Code)
	}
	resp := decodeResponse(t, rr)
	if resp.Error == "" {
		t.Error("expected error message for sqrt of negative")
	}
}

func TestCalculateHandler_SavesSuccessToHistory(t *testing.T) {
	hist := calculator.NewHistory()
	h := handlers.NewHandler(hist)

	post(t, h, `{"a":10,"b":2,"operation":"multiply"}`)

	entries := hist.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 history entry, got %d", len(entries))
	}
	if entries[0].Result != 20 {
		t.Errorf("expected history result 20, got %v", entries[0].Result)
	}
}

func TestCalculateHandler_DoesNotSaveErrorToHistory(t *testing.T) {
	hist := calculator.NewHistory()
	h := handlers.NewHandler(hist)

	post(t, h, `{"a":1,"b":0,"operation":"divide"}`)

	if len(hist.Entries()) != 0 {
		t.Error("failed calculations should not be saved to history")
	}
}

func TestHistoryHandler_Empty(t *testing.T) {
	h := handlers.NewHandler(calculator.NewHistory())
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	rr := httptest.NewRecorder()
	h.History(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var entries []calculator.HistoryEntry
	if err := json.NewDecoder(rr.Body).Decode(&entries); err != nil {
		t.Fatalf("failed to decode history: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty array, got %d entries", len(entries))
	}
}

func TestHistoryHandler_ReturnsEntries(t *testing.T) {
	hist := calculator.NewHistory()
	hist.Save(calculator.HistoryEntry{A: 2, B: 3, Operation: "add", Result: 5})
	hist.Save(calculator.HistoryEntry{A: 9, Operation: "sqrt", Result: 3})

	h := handlers.NewHandler(hist)
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	rr := httptest.NewRecorder()
	h.History(rr, req)

	var entries []calculator.HistoryEntry
	json.NewDecoder(rr.Body).Decode(&entries)
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestCalculateHandler_ContentTypeHeader(t *testing.T) {
	h := handlers.NewHandler(calculator.NewHistory())
	rr := post(t, h, `{"a":1,"b":1,"operation":"add"}`)
	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
}
