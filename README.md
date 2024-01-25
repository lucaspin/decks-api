An HTTP API for managing decks of cards.

- [Running the server](#running-the-server)
  - [With Docker](#with-docker)
  - [Without docker](#without-docker)
- [Running tests](#running-tests)
- [Storage implementations](#storage-implementations)
- [API](#api)
  - [Authentication](#authentication)
  - [Creating a deck](#creating-a-deck)
    - [Parameters](#parameters)
    - [Responses](#responses)
    - [Example - create a default deck (unshuffled, all cards)](#example---create-a-default-deck-unshuffled-all-cards)
    - [Example - create a shuffled deck (all cards)](#example---create-a-shuffled-deck-all-cards)
    - [Example - create an unshuffled deck with specific cards](#example---create-an-unshuffled-deck-with-specific-cards)
    - [Example - create a shuffled deck with specific cards](#example---create-a-shuffled-deck-with-specific-cards)
  - [Opening a deck](#opening-a-deck)
    - [Params](#params)
    - [Responses](#responses-1)
  - [Drawing cards from a deck](#drawing-cards-from-a-deck)
    - [Params](#params-1)
    - [Responses](#responses-2)
    - [Example - draw single card from deck](#example---draw-single-card-from-deck)


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

Note: you'll need to have Go 1.21 installed on your machine.

## Running tests

Tests are run with the `make test` command.

## Storage implementations

The persistence of decks is done through the [Storage interface](./pkg/storage/storage.go). The current implementations available are:
- **In-memory**: the default one. Keeps all the decks in memory. All the decks are lost if the server is shutdown.
- **Redis**: a Redis one. Note that this implementation has a few caveats currently, explained in [here](./pkg/storage/redis_storage.go). To use it, set the `DECK_STORAGE_TYPE` to `redis`.

## API

### Authentication

There was no requirement about authentication on the task description, so I decided not to implement it. The API is currently behind no authentication. However, I did register a [auth middleware](./pkg/api/auth.go), so if authentication is needed, that would be a good place to put it.

### Creating a deck

```
GET /api/v1alpha/decks
```

#### Parameters

- `shuffled` (optional) - determines if the cards in the deck will be shuffled or not. Default: false.
- `cards` (optional) - comma-separated list of card codes to include in the deck. If this is not specified, a deck with all 52 cards is created.

#### Responses

<b>201 Created</b>

```json
{
  "deck_id": "289970dd-32b0-4c88-a4c0-d2b2d1fbc53c",
  "shuffled": false,
  "remaining": 52
}
```

<b>400 Bad Request</b>

If the card codes specified in the `cards` parameter contains an invalid code, a 400 is returned.

#### Example - create a default deck (unshuffled, all cards)

```
curl -X POST http://localhost:4000/api/v1alpha/decks
```

#### Example - create a shuffled deck (all cards)

```
curl -X POST http://localhost:4000/api/v1alpha/decks?shuffled=true
```

#### Example - create an unshuffled deck with specific cards

```
curl -X POST http://localhost:4000/api/v1alpha/decks?cards=AH,2C,3D,KS
```

#### Example - create a shuffled deck with specific cards

```
curl -X POST http://localhost:4000/api/v1alpha/decks?cards=AH,2C,3D,KS&shuffled=true
```

### Opening a deck

```
GET /api/v1alpha/decks/:deck_id
```

#### Params

- `deck_id` (**required**) - the ID of the deck to open

#### Responses

<b>200 OK</b>

```json
{
  "deck_id": "bbf72234-b1a7-4671-aa47-1d75a99476a7",
  "shuffled": true,
  "remaining": 4,
  "cards": [
    {
      "Value": "KING",
      "Suit": "SPADES",
      "Code": "KS"
    },
    {
      "Value": "2",
      "Suit": "CLUBS",
      "Code": "2C"
    },
    {
      "Value": "ACE",
      "Suit": "HEARTS",
      "Code": "AH"
    },
    {
      "Value": "3",
      "Suit": "DIAMONDS",
      "Code": "3D"
    }
  ]
}
```

<b>400 Bad Request</b>

If the `deck_id` specified is not a valid UUID, 400 is returned.

<b>404 Not Found</b>

If the `deck_id` specified does not exist, 404 is returned.

### Drawing cards from a deck

```
POST /api/v1alpha/decks/:deck_id/draw
```

#### Params

- `deck_id` (**required**) - the ID of the deck to draw cards from.
- `count` (**required**) - how many cards to draw from the deck. This must be a positive integer. If this number is bigger than the current number of cards in the deck, all the cards in the deck are returned.

#### Responses

<b>200 OK</b>

```json
{
  "cards": [
    {
      "Value": "KING",
      "Suit": "SPADES",
      "Code": "KS"
    },
    {
      "Value": "2",
      "Suit": "CLUBS",
      "Code": "2C"
    }
  ]
}
```

<b>400 Bad Request</b>

A 400 status code is returned when:
- The `deck_id` specified is not a valid UUID.
- The `count` parameter is not specified, or it is not a valid positive integer.
- The deck is already empty.

<b>404 Not Found</b>

If the `deck_id` specified does not exist, 404 is returned.

#### Example - draw single card from deck

```
curl -X POST http://localhost:4000/api/v1alpha/decks/{deck_id}/draw?count=1
```
