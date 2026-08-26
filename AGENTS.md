# AGENTS.md

HTTP API for managing decks of cards.

## Quick facts

- Language: Go `1.21.4` (see `go.mod`)
- Module path: `github.com/lucaspin/decks-api`
- HTTP framework: [gorilla/mux](https://github.com/gorilla/mux)
- Authentication: none implemented yet
- Full API reference (routes, params, example curl commands, response shapes): [README.md](./README.md)

## Project layout

```
main.go          - entrypoint, wires up storage and starts the HTTP server
pkg/api/         - HTTP layer: router, handlers, request/response types, auth middleware
pkg/cards/       - card and deck domain logic (ranks, suits, card generation)
pkg/storage/     - Storage interface + in-memory and Redis implementations
```

## Build & run

```bash
make build
./build/server
```

Set `API_PORT` to run on a different port (defaults to `4000`):

```bash
API_PORT=8012 ./build/server
```

### With Docker

Requires Docker and Docker Compose.

```bash
make server.start   # build and start the server (also starts a Redis container)
make server.stop    # stop the server
make server.logs    # tail the server logs
```

## Tests

```bash
make test
```

This runs `gotestsum` inside the `app` Docker container (see `docker-compose.yml` / `Dockerfile.dev`, Go 1.21 image). Tests live alongside the code they cover as `_test.go` files:

- `pkg/api/server_test.go`
- `pkg/cards/card_generator_test.go`
- `pkg/cards/card_test.go`
- `pkg/storage/storage_test.go`

[testify](https://github.com/stretchr/testify) is used for assertions.

## Storage abstraction

Persistence is done through the `storage.Storage` interface (`pkg/storage/storage.go`), with three methods: `Create`, `Get`, `Draw`. Two implementations are available, selected via the `DECK_STORAGE_TYPE` env var:

- **in-memory** (`pkg/storage/in_memory_storage.go`) - the default. Decks are lost on restart.
- **redis** (`pkg/storage/redis_storage.go`) - set `DECK_STORAGE_TYPE=redis`. Has known caveats documented in that file - read it before relying on this backend.

Errors from the storage layer are sentinel values (`storage.ErrDeckNotFound`, `storage.ErrEmptyDeck`) checked with `errors.Is` and translated into HTTP status codes in `pkg/api/server.go`.

## API surface

All routes below are handled in `pkg/api/server.go`. See [README.md](./README.md) for full request/response documentation - avoid duplicating it here.

- `POST /api/v1alpha/decks` - create a deck
- `GET /api/v1alpha/decks/{deck_id}` - open a deck
- `POST /api/v1alpha/decks/{deck_id}/draw` - draw cards from a deck
- `GET /` - health check

## Conventions & style

- Standard Go formatting (`gofmt`).
- Storage errors are sentinel values, checked with `errors.Is` and mapped to HTTP status codes in the handlers (`pkg/api/server.go`).
- JSON responses are built with `respondWithJSON` and the response constructors in `pkg/api/responses.go` - follow this pattern for new endpoints rather than marshaling ad hoc structs in handlers.
- `pkg/api/auth.go` contains a no-op `authMiddleware` placeholder - this is the intended hook if authentication is ever added.

## Contribution guidance for agents

- Run `gofmt -l .` and `make test` before committing.
- Keep `README.md` in sync with any change to API behavior, request/response shapes, or storage backends.
- Avoid introducing new dependencies without checking `go.mod`/`go.sum`; run `go mod tidy` if you do.
