package calculator

import (
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	cases := []struct{ a, b, want float64 }{
		{3, 4, 7},
		{-1, 1, 0},
		{0, 0, 0},
		{1.5, 2.5, 4},
	}
	for _, tc := range cases {
		got, err := Add{}.Execute(tc.a, tc.b)
		if err != nil || got != tc.want {
			t.Errorf("Add(%v, %v) = %v, %v; want %v", tc.a, tc.b, got, err, tc.want)
		}
	}
}

func TestSubtract(t *testing.T) {
	cases := []struct{ a, b, want float64 }{
		{10, 3, 7},
		{0, 5, -5},
		{-2, -3, 1},
	}
	for _, tc := range cases {
		got, err := Subtract{}.Execute(tc.a, tc.b)
		if err != nil || got != tc.want {
			t.Errorf("Subtract(%v, %v) = %v, %v; want %v", tc.a, tc.b, got, err, tc.want)
		}
	}
}

func TestMultiply(t *testing.T) {
	cases := []struct{ a, b, want float64 }{
		{3, 4, 12},
		{-2, 5, -10},
		{0, 100, 0},
	}
	for _, tc := range cases {
		got, err := Multiply{}.Execute(tc.a, tc.b)
		if err != nil || got != tc.want {
			t.Errorf("Multiply(%v, %v) = %v, %v; want %v", tc.a, tc.b, got, err, tc.want)
		}
	}
}

func TestDivide(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		got, err := Divide{}.Execute(10, 2)
		if err != nil || got != 5 {
			t.Errorf("expected 5, got %v (err: %v)", got, err)
		}
	})
	t.Run("division by zero", func(t *testing.T) {
		_, err := Divide{}.Execute(5, 0)
		if err == nil {
			t.Error("expected error for division by zero, got nil")
		}
	})
}

func TestSqrt(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		got, err := Sqrt{}.Execute(9, 0)
		if err != nil || got != 3 {
			t.Errorf("expected 3, got %v (err: %v)", got, err)
		}
	})
	t.Run("negative input", func(t *testing.T) {
		_, err := Sqrt{}.Execute(-1, 0)
		if err == nil {
			t.Error("expected error for sqrt of negative number, got nil")
		}
	})
	t.Run("zero", func(t *testing.T) {
		got, err := Sqrt{}.Execute(0, 0)
		if err != nil || got != 0 {
			t.Errorf("expected 0, got %v (err: %v)", got, err)
		}
	})
}

func TestPercentage(t *testing.T) {
	cases := []struct{ a, want float64 }{
		{25, 0.25},
		{100, 1},
		{0, 0},
		{50, 0.5},
	}
	for _, tc := range cases {
		got, err := Percentage{}.Execute(tc.a, 0)
		if err != nil || math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("Percentage(%v) = %v, %v; want %v", tc.a, got, err, tc.want)
		}
	}
}

func TestPower(t *testing.T) {
	cases := []struct{ a, b, want float64 }{
		{2, 10, 1024},
		{3, 3, 27},
		{5, 0, 1},
		{0, 0, 1},
		{2, -1, 0.5},
	}
	for _, tc := range cases {
		got, err := Power{}.Execute(tc.a, tc.b)
		if err != nil || math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("Power(%v, %v) = %v, %v; want %v", tc.a, tc.b, got, err, tc.want)
		}
	}
}

func TestGetOperation(t *testing.T) {
	valid := []string{"add", "subtract", "multiply", "divide", "sqrt", "percentage", "power"}
	for _, name := range valid {
		if _, ok := GetOperation(name); !ok {
			t.Errorf("expected operation %q to exist in registry", name)
		}
	}
	if _, ok := GetOperation("unknown"); ok {
		t.Error("expected unknown operation to not exist in registry")
	}
}

func TestCalculatorCompute(t *testing.T) {
	op, _ := GetOperation("add")
	calc := NewCalculator(op)
	result, err := calc.Compute(6, 4)
	if err != nil || result != 10 {
		t.Errorf("expected 10, got %v (err: %v)", result, err)
	}
}

func TestHistory(t *testing.T) {
	h := NewHistory()
	if len(h.Entries()) != 0 {
		t.Error("expected empty history")
	}

	h.Save(HistoryEntry{A: 1, B: 2, Operation: "add", Result: 3})
	h.Save(HistoryEntry{A: 9, Operation: "sqrt", Result: 3})

	entries := h.Entries()
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Result != 3 || entries[0].Operation != "add" {
		t.Errorf("unexpected first entry: %+v", entries[0])
	}
}
