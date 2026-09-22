An HTTP API for managing decks of cards.

- [Running the server](#running-the-server)
  - [With Docker](#with-docker)
  - [Without docker](#without-docker)
- [Running tests](#running-tests)
- [Storage implementations](#storage-implementations)
- [API](#api)


## Running the server

### With Docker

If you have Docker and Docker Compose on your machine, you can start the server with the `make server.start` command. This command will start the server inside of a Docker container, and the server will be available on `localhost:4000`.

You can also:
- Stop the server with `make server.stop`
- Inspect its logs with `make server.logs`
- Make the server start on a different by specifying the `API_PORT` environment variable in `docker-compose.yml`

### Without docker

If you don't want to run the server with Docker, you can build it and run it directly:

```bash
make build
./build/server
```

You can also specify a different port for the server to start:

```bash
API_PORT=8012 ./build/server
```

Note: you'll need to have Go 1.27 installed on your machine.

## Running tests

Tests are run with the `make test` command.

## Storage implementations

The persistence of decks is done through the [Storage interface](./pkg/storage/storage.go). The current implementations available are:
- **In-memory**: the default one. Keeps all the decks in memory. All the decks are lost if the server is shutdown.
- **Redis**: a Redis one. Note that this implementation has a few caveats currently, explained in [here](./pkg/storage/redis_storage.go). To use it, set the `DECK_STORAGE_TYPE` to `redis`.

## API

See [docs/API.md](./docs/API.md) for the full API reference: authentication and all endpoints (creating, opening, drawing from, and deleting decks), including parameters, response shapes, and example curl commands.
