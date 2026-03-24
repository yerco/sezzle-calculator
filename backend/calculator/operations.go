package calculator

import (
	"errors"
	"math"
)

// Operation is the Strategy interface — each arithmetic op is a concrete strategy.
type Operation interface {
	Execute(a, b float64) (float64, error)
}

// --- Concrete strategies ---

type Add struct{}

func (Add) Execute(a, b float64) (float64, error) { return a + b, nil }

type Subtract struct{}

func (Subtract) Execute(a, b float64) (float64, error) { return a - b, nil }

type Multiply struct{}

func (Multiply) Execute(a, b float64) (float64, error) { return a * b, nil }

type Divide struct{}

func (Divide) Execute(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// Sqrt is a unary operation; b is ignored.
type Sqrt struct{}

func (Sqrt) Execute(a, b float64) (float64, error) {
	if a < 0 {
		return 0, errors.New("square root of negative number")
	}
	return math.Sqrt(a), nil
}

// Percentage converts a to its decimal form (25 → 0.25); b is ignored.
type Percentage struct{}

func (Percentage) Execute(a, b float64) (float64, error) { return a / 100, nil }

// Power raises a to the power of b (aᵇ).
type Power struct{}

func (Power) Execute(a, b float64) (float64, error) { return math.Pow(a, b), nil }

// --- Calculator context (Strategy pattern) ---

// Calculator holds and executes a chosen Operation strategy.
type Calculator struct {
	op Operation
}

func NewCalculator(op Operation) *Calculator {
	return &Calculator{op: op}
}

func (c *Calculator) Compute(a, b float64) (float64, error) {
	return c.op.Execute(a, b)
}

// --- Registry ---

// registry maps operation names to their Strategy implementations.
var registry = map[string]Operation{
	"add":        Add{},
	"subtract":   Subtract{},
	"multiply":   Multiply{},
	"divide":     Divide{},
	"sqrt":       Sqrt{},
	"percentage": Percentage{},
	"power":      Power{},
}

// GetOperation returns the Operation for the given name, or false if unknown.
func GetOperation(name string) (Operation, bool) {
	op, ok := registry[name]
	return op, ok
}
