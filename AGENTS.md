# AGENTS.md

Instructions for AI coding agents working on this repository.

## Project overview

`decks-api` is an HTTP API for managing decks of cards, written in Go
(module `github.com/lucaspin/decks-api`, Go 1.21.4). It uses `gorilla/mux`
for routing and `gorilla/handlers` for middleware, and stores decks through
a pluggable `Storage` interface (in-memory or Redis).

## Setup / running locally

### With Docker (preferred)

```bash
make server.start   # builds and starts the server + Redis via docker-compose
make server.stop    # stops everything
make server.logs    # tails the app container logs
```

The server listens on `localhost:4000` by default.

### Without Docker

```bash
make build       # go build -o build/server main.go
./build/server
```

Requires Go 1.21 installed locally. Override the port with `API_PORT`
(e.g. `API_PORT=8012 ./build/server`).

## Testing

```bash
make test
```

This runs `gotestsum` inside the `app` container (via `docker-compose run`)
against `./...`. Tests live next to the code they cover (`*_test.go`), use
`github.com/stretchr/testify/require` for assertions, and follow the
`Test__<Thing>` naming convention with table-driven cases exercised via
`t.Run` subtests. Follow this style for any new tests.

## Code structure

- `main.go` - entrypoint; wires up storage via `storage.NewStorage()`,
  reads `API_PORT`, and starts the server via `api.NewServer(...).Serve(...)`.
- `pkg/api` - HTTP layer: routes under `/api/v1alpha` (`server.go`), response
  helpers (`responses.go`), an auth middleware stub (`auth.go`), and
  `server_test.go` with `httptest`-based handler tests.
- `pkg/cards` - card and deck generation logic (parsing card codes,
  building/shuffling decks) plus its tests.
- `pkg/storage` - the `Storage` interface (`storage.go`), an in-memory
  implementation (`in_memory_storage.go`), a Redis implementation
  (`redis_storage.go`), and shared `storage_test.go` covering both.

Add new HTTP endpoints in `pkg/api`, card/deck logic in `pkg/cards`, and
persistence logic in `pkg/storage`.

## Storage backends

The `Storage` interface is defined in `pkg/storage/storage.go`. The backend
is selected via the `DECK_STORAGE_TYPE` environment variable:

- unset / anything else -> in-memory storage (default; decks are lost on
  restart).
- `redis` -> Redis-backed storage. See the caveats documented at the top of
  `pkg/storage/redis_storage.go` (operations are not atomic, so it is not
  safe for concurrent requests against the same deck without extra
  locking).

Any change to storage behavior must be reflected in both implementations so
they keep satisfying the same interface/contract, and `storage_test.go`
exercises both.

## API surface

Full endpoint documentation (routes, params, responses, curl examples, and
the authentication note) lives in `README.md` under the "API" section.
Refer to it instead of duplicating it here, and keep it up to date whenever
API behavior changes.

## Conventions for agents

- Before committing, make sure the code builds and is formatted:
  `gofmt -l .` and `go vet ./...` (or `go build ./...`) should be clean.
- Add or update tests for any behavior change, following the existing
  table-driven / `t.Run` / `Test__Name` style with `testify/require`.
- Don't break the `Storage` interface contract; keep the in-memory and
  Redis implementations in sync.
- Keep `README.md` in sync with any API changes, since this file
  intentionally defers to it for API documentation.
- Keep changes focused and minimal; don't commit the `build/` output
  produced by `make build` (already covered by `.gitignore`).
