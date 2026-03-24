# Sezzle Calculator

A full-stack calculator application built as a take-home assignment for the Engineering Manager role at Sezzle.

The goal was not just to build a working calculator — it was to build one with clean interfaces, deliberate patterns, testable code, and zero unnecessary complexity.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | React 19, TypeScript, Vite 5 |
| Backend | Go 1.22, `net/http` (stdlib only) |
| Testing | Go `testing` package, Vitest 2, React Testing Library |
| Containers | Docker + Docker Compose |

No third-party backend frameworks. The Go standard library in 1.22 has everything needed — adding Gin or Echo for a project this size would be over-engineering.

---

## Project Structure

```
sezzle-calculator/
├── backend/
│   ├── main.go                    # Server setup, routing, CORS middleware
│   ├── go.mod
│   ├── Dockerfile
│   ├── calculator/
│   │   ├── operations.go          # Strategy pattern: all arithmetic operations + registry
│   │   ├── history.go             # Memento pattern: HistoryEntry + thread-safe History
│   │   └── operations_test.go     # Unit tests for all operations, registry, and history
│   ├── handlers/
│   │   └── calculator.go          # HTTP handlers with injected dependencies
│   ├── middleware/
│   │   └── validation.go          # Input validation + context injection
│   └── models/
│       └── request.go             # CalculationRequest / CalculationResponse structs
├── frontend/
│   ├── src/
│   │   ├── api/
│   │   │   └── calculator.ts      # fetch wrappers for /calculate and /history
│   │   ├── components/
│   │   │   ├── Calculator.tsx     # main calculator UI + button grid
│   │   │   ├── Calculator.test.tsx
│   │   │   └── History.tsx        # history panel
│   │   ├── hooks/
│   │   │   ├── useCalculator.ts   # all calculator state and API logic
│   │   │   └── useCalculator.test.ts
│   │   └── App.tsx
│   ├── package.json
│   └── Dockerfile
├── docker-compose.yml
└── README.md
```

---

## Getting Started

### Prerequisites

- Go 1.22+
- Node.js 18+ (20+ recommended; required for Docker builds)
- Docker + Docker Compose (optional)

### Run with Docker Compose (recommended)

```bash
docker-compose up --build
```

| Service | URL |
|---|---|
| Frontend | http://localhost:5173 |
| Backend API | http://localhost:8080 |

### Run Locally

**Backend:**
```bash
cd backend
go run main.go
# API available at http://localhost:8080
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
# UI available at http://localhost:5173
```

---

## API Reference

### `POST /calculate`

Performs a calculation and saves a snapshot to the in-memory history.

**Request:**
```json
{
  "a": 10,
  "b": 5,
  "operation": "divide"
}
```

**Success response `200`:**
```json
{ "result": 2 }
```

**Supported operations:**

| Operation | Description | Example |
|---|---|---|
| `add` | `a + b` | `3 + 4 = 7` |
| `subtract` | `a - b` | `10 - 3 = 7` |
| `multiply` | `a × b` | `3 × 4 = 12` |
| `divide` | `a ÷ b` | `10 ÷ 2 = 5` |
| `power` | `aᵇ` | `2 ^ 10 = 1024` |
| `sqrt` | `√a` (b ignored) | `√9 = 3` |
| `percentage` | `a ÷ 100` (b ignored) | `25% = 0.25` |

**Error responses:**

| Scenario | Status | Body |
|---|---|---|
| Invalid JSON | `400` | `{ "error": "invalid request body" }` |
| Unknown operation | `400` | `{ "error": "unsupported operation: foo" }` |
| Division by zero | `422` | `{ "error": "division by zero" }` |
| Sqrt of negative | `422` | `{ "error": "square root of negative number" }` |

---

### `GET /history`

Returns all successful calculations from the current session, in insertion order.

**Response `200`:**
```json
[
  { "a": 10, "b": 5, "operation": "add", "result": 15, "timestamp": "2026-03-23T14:00:00Z" },
  { "a": 9,  "b": 0, "operation": "sqrt", "result": 3,  "timestamp": "2026-03-23T14:00:05Z" }
]
```

Always returns a JSON array — never `null`. Failed calculations are not recorded.

---

## Design Decisions

### Strategy Pattern — arithmetic operations

Each operation (`Add`, `Subtract`, `Multiply`, `Divide`, `Power`, `Sqrt`, `Percentage`) is a struct that satisfies the `Operation` interface:

```go
type Operation interface {
    Execute(a, b float64) (float64, error)
}
```

The `Calculator` context holds a reference to whichever strategy is selected at runtime. A `registry` map decouples operation selection from the handler — adding a new operation means adding one struct and one registry entry, nothing else changes. Open/closed in practice, not just on paper.

### Memento Pattern — calculation history

`HistoryEntry` is an immutable snapshot of a single successful calculation. The `History` struct accumulates entries behind a `sync.RWMutex` — safe for concurrent requests. A `HistoryEntry` is only created after a successful `Compute` call, so failed calculations never appear in the history.

### Middleware for input validation

Validation is a cross-cutting concern and doesn't belong in the handler. The `ValidateCalculation` middleware decodes the request body, rejects unknown operations with a clear error, then passes the validated struct forward via `context.WithValue`. Handlers never parse raw JSON — they receive a typed `CalculationRequest`. This keeps handlers readable and validation independently testable.

### Dependency injection over global state

`Handler` is a struct that receives `*calculator.History` at construction time. `main.go` wires everything together. Nothing is global. This makes the handlers testable in isolation — pass in a fresh `History`, assert on its state.

### Error semantics

`400 Bad Request` for things the client sent wrong (malformed JSON, unknown operation). `422 Unprocessable Entity` for inputs that are structurally valid but semantically impossible (division by zero, sqrt of a negative). The distinction matters — a client that gets a `400` should fix its request; a `422` might be a valid user action that needs a friendly UI message.

### Docker — two-stage builds

Both Dockerfiles use multi-stage builds: the backend compiles to a static binary copied into a bare `alpine` image (~10 MB), the frontend builds with Node then serves the static output from `nginx:alpine`. No source code or build tooling ships in the final images.

---

## Running Tests

```bash
# Backend — all packages
cd backend && go test ./... -v

# Backend — with coverage
cd backend && go test ./... -cover

# Frontend
cd frontend && npm install && npm test
cd frontend && npm run test:coverage
```

**Backend** — 31 tests across 3 packages, all passing:

| Package | Tests | Coverage |
|---|---|---|
| `calculator` | 13 | 100% |
| `handlers` | 9 | 79% |
| `middleware` | 9 | 100% |

`main.go` (`cors`, `main`) is excluded from unit tests by design — covered by the Docker integration path. Full report: `backend/coverage.out` / `backend/coverage.txt`.

**Frontend** — 36 tests across 3 suites, all passing:

| Suite | Tests | Coverage |
|---|---|---|
| `api/calculator.test.ts` | 8 | 100% |
| `hooks/useCalculator.test.ts` | 18 | 99% |
| `components/Calculator.test.tsx` | 9 | 100% (Calculator), 43% (History) |

`History.tsx` formatting helpers and `App.tsx`/`main.tsx` are not unit-tested — they are pure rendering with no logic.

---

## What I Would Add Given More Time

These are deliberate omissions, not oversights:

- **Persistence** — swap the in-memory `History` for a storage interface with a DB-backed implementation. The injection pattern already supports this without touching handlers.
- **History pagination** — `GET /history?limit=20&offset=0`.
- **Auth** — session-scoped history so each user sees only their own calculations.
- **Structured logging** — `slog` (stdlib, Go 1.21+) for request logging and observability.
- **CI pipeline** — GitHub Actions running `go test ./...` and `npm run test` on every push.

---

## AI Tooling

This project was built using Claude (Anthropic) as a coding assistant, in line with the assignment's permission to use AI tools.

Claude was used to accelerate Go implementation and generate test scaffolding. Architecture decisions, pattern selection, error semantics, and code review were mine throughout.

I treat AI the same way I'd treat a fast junior engineer: useful for execution speed, needs direction on trade-offs, and everything it produces gets reviewed before it ships.

---

## Author

**Yerco Jorquera C.**
