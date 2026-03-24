package calculator

import (
	"sync"
	"time"
)

// HistoryEntry is the Memento — an immutable snapshot of a single calculation.
type HistoryEntry struct {
	A         float64   `json:"a"`
	B         float64   `json:"b"`
	Operation string    `json:"operation"`
	Result    float64   `json:"result"`
	Timestamp time.Time `json:"timestamp"`
}

// History is the Caretaker — it stores and retrieves HistoryEntry mementos.
// The Calculator handler acts as the Originator, creating entries after each
// successful computation.
type History struct {
	mu      sync.RWMutex
	entries []HistoryEntry
}

func NewHistory() *History {
	return &History{}
}

// Save appends a new memento to the history.
func (h *History) Save(entry HistoryEntry) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = append(h.entries, entry)
}

// Entries returns a copy of all stored mementos.
func (h *History) Entries() []HistoryEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	result := make([]HistoryEntry, len(h.entries))
	copy(result, h.entries)
	return result
}
