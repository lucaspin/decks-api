# AGENTS.md

Guidance for coding agents (and humans) working in this repository.

## Project overview

- **What it is**: an HTTP API for managing decks of cards. See [README.md](./README.md) for the full API reference (endpoints, params, example requests/responses).
- **Module**: `github.com/lucaspin/decks-api`, Go `1.21.4` (see `go.mod`).
- **Entry point**: `main.go` — wires up a `storage.Storage` implementation and an `api.Server`, then starts listening. Port is controlled by the `API_PORT` env var (default `4000`).

## Repo layout

- `main.go` — process entrypoint.
- `pkg/api` — HTTP server (`gorilla/mux` router), handlers, JSON response types, auth middleware stub.
  - `server.go` — `Server` type, `InitRouter`, handlers.
  - `responses.go` — response struct types + `respondWithJSON`.
  - `auth.go` — auth middleware stub (currently unused/no-op, see Caveats).
  - `server_test.go` — handler tests.
- `pkg/cards` — card and deck generation logic (`card.go`, `card_generator.go`) with matching `_test.go` files.
- `pkg/storage` — `Storage` interface (`storage.go`) with two implementations:
  - `in_memory_storage.go` — default, keeps decks in memory (lost on restart).
  - `redis_storage.go` — Redis-backed, selected via `DECK_STORAGE_TYPE=redis`. Has known atomicity caveats documented at the top of that file — read it before relying on Redis behavior.
  - `storage_test.go` — shared storage tests.
- `Dockerfile.dev`, `docker-compose.yml` — dev/runtime containers (`app` service + `redis` service).
- `Makefile` — canonical build/run/test commands (see below).

## Build, run, test

Use the Makefile targets; don't invent new commands.

- **Build**: `make build` — builds `./build/server`.
- **Run without Docker**: `make build && ./build/server` (requires Go 1.21 locally). Override the port with `API_PORT=8012 ./build/server`.
- **Run with Docker**: `make server.start` (builds + starts `app` and `redis` containers, server on `localhost:4000`), `make server.stop`, `make server.logs`.
- **Test**: `make test` runs `docker-compose run --rm app gotestsum --format short-verbose --packages="./..." -- -p 1`. This requires Docker/Docker Compose.
  - If Docker is not available in your environment, fall back to `go test ./...` directly. Be aware this may differ from CI: the Redis-backed storage tests need a reachable Redis (see `REDIS_HOST`/`REDIS_PORT`/`REDIS_USERNAME`/`REDIS_PASSWORD` env vars used in `docker-compose.yml`); without one, those tests will fail or need to be skipped/pointed at a local Redis instance.

## Conventions

- **Test naming**: `Test__Name` (double underscore), using `t.Run` subtests and `testify/require`. Follow the existing style in `pkg/api/server_test.go`, `pkg/cards/*_test.go`, `pkg/storage/storage_test.go`.
- **Handler error handling**: validate input first (→ 400), map storage sentinel errors (`storage.ErrDeckNotFound`, `storage.ErrEmptyDeck`) to the appropriate status codes, and fall back to a logged 500 for anything unexpected.
- **Storage backends**: must implement `pkg/storage.Storage` (`Create`/`Get`/`Draw`). Register new implementations as a case in `storage.NewStorage()`, keyed off `DECK_STORAGE_TYPE`.
- **Routing**: routes are versioned under `/api/v1alpha` and registered in `InitRouter` in `pkg/api/server.go`. Add new routes there.
- **Responses**: use small, dedicated struct types in `pkg/api/responses.go` and send them via `respondWithJSON`.
- **Formatting**: idiomatic Go, `gofmt`-clean; run `go vet ./...` before committing. There is no linter config in the repo currently.

## Known caveats

- **No authentication**: the API has no auth. `pkg/api/auth.go` contains a middleware stub intended as the extension point if auth is ever added.
- **Redis storage caveats**: `pkg/storage/redis_storage.go` documents that its operations are not atomic, so it isn't safe for concurrent requests against the same deck. Read the comment at the top of that file before depending on Redis storage semantics.

## Quick pointers for common tasks

- **Add an endpoint**: add a handler in `pkg/api/server.go`, register the route in `InitRouter`, add a response type in `pkg/api/responses.go`, add tests in `pkg/api/server_test.go`.
- **Add a storage backend**: implement `pkg/storage.Storage`, add a case in `storage.NewStorage()`, add tests mirroring `pkg/storage/storage_test.go`.
- **Add/modify card logic**: `pkg/cards/card.go` and `pkg/cards/card_generator.go`, with corresponding `_test.go` files.
