package models

// CalculationRequest represents the API input
type CalculationRequest struct {
	A         float64 `json:"a"`
	B         float64 `json:"b"`
	Operation string  `json:"operation"`
}

// CalculationResponse represents the API output.
// Result is a pointer so that zero values (e.g. 0+0=0) are not omitted by omitempty.
type CalculationResponse struct {
	Result *float64 `json:"result,omitempty"`
	Error  string   `json:"error,omitempty"`
}
