# AGENTS.md

Guidance for coding agents working in this repository.

## Project overview

`decks-api` (Go module `github.com/lucaspin/decks-api`, Go 1.21.4) is an HTTP
API for managing decks of cards: creating decks, opening them, and drawing
cards from them. Full API documentation with request/response examples lives
in [README.md](./README.md) — this file focuses on codebase orientation and
contribution conventions, not on duplicating that documentation.

## Repo layout

- `main.go` — entrypoint. Initializes storage via `storage.NewStorage()`,
  wires it into `api.NewServer`, and starts the server on the port from
  `API_PORT` (default `4000`).
- `pkg/api` — HTTP layer.
  - `server.go` — `Server` type, router setup (`InitRouter`), and handlers
    (`CreateDeck`, `OpenDeck`, `DrawCards`, `HealthCheck`).
  - `responses.go` — response shaping helpers (e.g. `respondWithJSON`).
  - `auth.go` — no-op auth middleware placeholder.
  - `server_test.go` — handler tests.
- `pkg/storage` — persistence layer.
  - `storage.go` — `Storage` interface (`Create`/`Get`/`Draw`), `Deck` type,
    and `NewStorage()` factory that picks an implementation based on
    `DECK_STORAGE_TYPE`.
  - `in_memory_storage.go` — default, in-process implementation. Decks are
    lost on restart.
  - `redis_storage.go` — Redis-backed implementation (opt-in). Documented in
    its own comment as **not atomic/safe for concurrent access to the same
    deck** — see the note below before relying on it.
  - `storage_test.go` — tests for both implementations.
- `pkg/cards` — domain types and generation logic.
  - `card.go` — `Card`, `CardSuit`, `CardRank` types.
  - `card_generator.go` — deck/card generation, shuffling, and card-code
    parsing (e.g. `AH`, `2C`, `KS`).
  - `card_test.go`, `card_generator_test.go` — tests.
- `docker-compose.yml`, `Dockerfile.dev` — local dev environment (`app` +
  `redis` containers).
- `Makefile` — build, test, and server start/stop/logs targets.

## Build, run, and test

- Build: `make build` → produces `./build/server`. Run it directly with
  `./build/server` (Go 1.21 required locally).
- Run with Docker: `make server.start` (serves on `localhost:4000`),
  `make server.stop`, `make server.logs`.
- Tests: `make test` — runs
  `docker-compose run --rm app gotestsum --format short-verbose --packages="./..." -- -p 1`.
  This requires Docker/Docker Compose, since tests run inside the `app`
  container with packages executed serially (`-p 1`). Run this before
  considering any change complete.

## Configuration / environment variables

- `API_PORT` — server port (default `4000`).
- `DECK_STORAGE_TYPE` — `"redis"` to use Redis storage; anything else (or
  unset) uses the in-memory default.
- `REDIS_HOST`, `REDIS_PORT`, `REDIS_USERNAME`, `REDIS_PASSWORD` — used by the
  Redis storage implementation (see `docker-compose.yml` for local values).

## Storage backends

Persistence lives in `pkg/storage`. All backends implement the same
`Storage` interface, and the active backend is chosen at startup.

- Interface (`storage.go`): `Storage` exposes `Create`, `Get`, and `Draw`.
  Shared types live here too — `Deck` (with `Remaining()`), and the sentinel
  errors `ErrDeckNotFound` and `ErrEmptyDeck` that handlers map to `404`/`400`.
- Selection (`storage.NewStorage()`): reads `DECK_STORAGE_TYPE`. `"redis"`
  selects the Redis backend; any other value (or unset) falls back to the
  in-memory backend.

Available backends:

- In-memory (`in_memory_storage.go`) — the default, built by
  `NewInMemoryStorage()`. Keeps decks in a `map` in process. Simple and
  dependency-free, but all decks are lost on restart and it is only intended
  for local use/tests.
- Redis (`redis_storage.go`) — opt-in, built by `NewRedisStorage(config)`.
  Reads connection settings from `REDIS_HOST`, `REDIS_PORT`, `REDIS_USERNAME`,
  and `REDIS_PASSWORD` when no config is passed. Each deck is stored as two
  keys: `decks:{deckID}:cards` (a Redis list of card codes, drawn via `LPOP`)
  and `decks:{deckID}:shuffled` (a flag). It is **not atomic and not safe for
  concurrent access to the same deck** — see the comment at the top of
  `redis_storage.go`; making it safe would require distributed locks or Lua
  scripts. Don't treat it as production-hardened without addressing that.

Adding a new backend:

1. Implement the `Storage` interface (`Create`/`Get`/`Draw`), returning
   `ErrDeckNotFound`/`ErrEmptyDeck` where appropriate so handlers keep their
   status-code behavior.
2. Wire it into `storage.NewStorage()` with a new `DECK_STORAGE_TYPE` value.
3. Add tests in `storage_test.go` covering the new implementation.

## API surface

- `POST /api/v1alpha/decks`
- `GET /api/v1alpha/decks/{deck_id}`
- `POST /api/v1alpha/decks/{deck_id}/draw`
- `GET /` (health check)

See [README.md](./README.md#api) for parameters, status codes, and example
requests/responses.

## Conventions for making changes

- New storage backends: see the [Storage backends](#storage-backends)
  section above for the interface, selection, and steps to add one.
- HTTP handlers live on `api.Server` in `pkg/api/server.go`. Register new
  routes in `InitRouter`, and follow existing status-code conventions:
  `400` for invalid input, `404` for not found, `500` for unknown/internal
  errors.
- Authentication is intentionally a no-op in `pkg/api/auth.go`. If auth is
  ever required, implement it there rather than adding ad hoc checks in
  handlers.
- Card code parsing/generation (e.g. `AH`, `2C`, `KS`) lives in
  `pkg/cards/card_generator.go` — reuse it instead of re-implementing parsing
  elsewhere.
- Keep `README.md` and this file in sync if commands, environment variables,
  or endpoints change.

## Testing expectations

- Each package has `_test.go` files using `testify`. Add or extend tests
  alongside code changes.
- Run `make test` before considering a change complete.

## Notes and caveats

- The Redis storage backend is explicitly non-atomic and not safe for
  concurrent access to the same deck — see the
  [Storage backends](#storage-backends) section for details.
- There is no authentication by design, not by omission — see
  `pkg/api/auth.go` and the README's Authentication section.
