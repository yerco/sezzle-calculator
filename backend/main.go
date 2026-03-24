package main

import (
	"fmt"
	"log"
	"net/http"

	"sezzle-calculator/calculator"
	"sezzle-calculator/handlers"
	"sezzle-calculator/middleware"
)

// cors wraps a handler to allow cross-origin requests from the frontend.
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	history := calculator.NewHistory()
	h := handlers.NewHandler(history)

	mux := http.NewServeMux()
	// Go 1.22 method-based routing.
	mux.HandleFunc("POST /calculate", middleware.ValidateCalculation(h.Calculate))
	mux.HandleFunc("GET /history", h.History)

	addr := ":8080"
	fmt.Printf("Sezzle Calculator API listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, cors(mux)))
}
