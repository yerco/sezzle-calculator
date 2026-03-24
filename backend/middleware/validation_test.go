package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"sezzle-calculator/middleware"
	"sezzle-calculator/models"
)

// sentinel is a handler that records whether it was reached.
func sentinel(reached *bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		*reached = true
	}
}

func runMiddleware(t *testing.T, body string) (rr *httptest.ResponseRecorder, reached bool) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/calculate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	middleware.ValidateCalculation(sentinel(&reached)).ServeHTTP(rr, req)
	return rr, reached
}

func TestValidateCalculation_ValidRequest(t *testing.T) {
	rr, reached := runMiddleware(t, `{"a":3,"b":4,"operation":"add"}`)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	if !reached {
		t.Error("expected next handler to be called")
	}
}

func TestValidateCalculation_SetsRequestInContext(t *testing.T) {
	var got models.CalculationRequest
	handler := func(w http.ResponseWriter, r *http.Request) {
		req, ok := middleware.GetRequest(r.Context())
		if !ok {
			t.Error("expected request in context")
		}
		got = req
	}
	req := httptest.NewRequest(http.MethodPost, "/calculate",
		strings.NewReader(`{"a":5,"b":2,"operation":"multiply"}`))
	req.Header.Set("Content-Type", "application/json")
	middleware.ValidateCalculation(handler).ServeHTTP(httptest.NewRecorder(), req)

	if got.A != 5 || got.B != 2 || got.Operation != "multiply" {
		t.Errorf("unexpected request in context: %+v", got)
	}
}

func TestValidateCalculation_InvalidJSON(t *testing.T) {
	rr, reached := runMiddleware(t, `not json`)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
	if reached {
		t.Error("next handler should not be called on invalid JSON")
	}
}

func TestValidateCalculation_UnknownFields(t *testing.T) {
	rr, reached := runMiddleware(t, `{"a":1,"b":2,"operation":"add","extra":"field"}`)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for unknown fields, got %d", rr.Code)
	}
	if reached {
		t.Error("next handler should not be called with unknown fields")
	}
}

func TestValidateCalculation_TrailingGarbage(t *testing.T) {
	rr, reached := runMiddleware(t, `{"a":1,"b":2,"operation":"add"}GARBAGE`)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for trailing garbage, got %d", rr.Code)
	}
	if reached {
		t.Error("next handler should not be called with trailing garbage")
	}
}

func TestValidateCalculation_UnknownOperation(t *testing.T) {
	rr, reached := runMiddleware(t, `{"a":1,"b":2,"operation":"modulo"}`)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
	if reached {
		t.Error("next handler should not be called with unknown operation")
	}
}

func TestValidateCalculation_UnaryWithNonZeroB(t *testing.T) {
	for _, op := range []string{"sqrt", "percentage"} {
		rr, reached := runMiddleware(t, `{"a":9,"b":5,"operation":"`+op+`"}`)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("%s: expected 400 for b!=0, got %d", op, rr.Code)
		}
		if reached {
			t.Errorf("%s: next handler should not be called when b!=0", op)
		}
	}
}

func TestValidateCalculation_UnaryWithZeroB(t *testing.T) {
	for _, op := range []string{"sqrt", "percentage"} {
		_, reached := runMiddleware(t, `{"a":9,"b":0,"operation":"`+op+`"}`)
		if !reached {
			t.Errorf("%s: next handler should be called when b=0", op)
		}
	}
}

func TestValidateCalculation_SetsContentTypeHeader(t *testing.T) {
	rr, _ := runMiddleware(t, `not json`)
	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
}
