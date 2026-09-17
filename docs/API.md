# API

All endpoints are served under `/api/v1alpha`. The base URL in the examples is `http://localhost:4000`.

## Authentication

The API is currently behind no authentication, as it was not part of the task. An [auth middleware](../pkg/api/auth.go) is registered as a placeholder for when it is needed.

## Creating a deck

```
POST /decks
```

Parameters:
- `shuffled` (optional) - whether the cards are shuffled. Default: `false`.
- `cards` (optional) - comma-separated card codes to include. Defaults to all 52 cards.

Responses:
- `201 Created` - `{ "deck_id": "...", "shuffled": false, "remaining": 52 }`
- `400 Bad Request` - a code in `cards` is invalid.

Examples:

```
curl -X POST http://localhost:4000/api/v1alpha/decks
curl -X POST "http://localhost:4000/api/v1alpha/decks?shuffled=true"
curl -X POST "http://localhost:4000/api/v1alpha/decks?cards=AH,2C,3D,KS"
curl -X POST "http://localhost:4000/api/v1alpha/decks?cards=AH,2C,3D,KS&shuffled=true"
```

## Opening a deck

```
GET /decks/:deck_id
```

Responses:
- `200 OK` - deck with `deck_id`, `shuffled`, `remaining`, and the full `cards` list.
- `400 Bad Request` - `deck_id` is not a valid UUID.
- `404 Not Found` - `deck_id` does not exist.

## Drawing cards from a deck

```
POST /decks/:deck_id/draw
```

Parameters:
- `count` (**required**) - positive integer number of cards to draw. If larger than the remaining count, all cards are returned.

Responses:
- `200 OK` - `{ "cards": [ ... ] }`
- `400 Bad Request` - invalid `deck_id`, missing/invalid `count`, or empty deck.
- `404 Not Found` - `deck_id` does not exist.

Example:

```
curl -X POST "http://localhost:4000/api/v1alpha/decks/{deck_id}/draw?count=1"
```

## Re-shuffling a deck

```
POST /decks/:deck_id/shuffle
```

Randomizes the order of the remaining cards. Already-drawn cards are not affected.

Responses:
- `200 OK` - `{ "deck_id": "...", "shuffled": true, "remaining": 52 }`
- `400 Bad Request` - `deck_id` is not a valid UUID.
- `404 Not Found` - `deck_id` does not exist.

Example:

```
curl -X POST http://localhost:4000/api/v1alpha/decks/{deck_id}/shuffle
```

## Deleting a deck

```
DELETE /decks/:deck_id
```

Responses:
- `204 No Content` - the deck was deleted.
- `400 Bad Request` - `deck_id` is not a valid UUID.
- `404 Not Found` - `deck_id` does not exist.

Example:

```
curl -X DELETE http://localhost:4000/api/v1alpha/decks/{deck_id}
```
