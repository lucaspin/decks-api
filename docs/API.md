# API

This document describes the HTTP API exposed by the server.

### Authentication

There was no requirement about authentication on the task description, so I decided not to implement it. The API is currently behind no authentication. However, I did register a [auth middleware](../pkg/api/auth.go), so if authentication is needed, that would be a good place to put it.

### Creating a deck

```
POST /api/v1alpha/decks
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

### Re-shuffling a deck

```
POST /api/v1alpha/decks/:deck_id/shuffle
```

Re-shuffles the cards currently remaining in the deck. Cards that have already been drawn are **not** brought back into the deck - only the remaining cards are re-shuffled.

#### Params

- `deck_id` (**required**) - the ID of the deck to re-shuffle.

#### Responses

<b>200 OK</b>

```json
{
  "deck_id": "289970dd-32b0-4c88-a4c0-d2b2d1fbc53c",
  "shuffled": true,
  "remaining": 52
}
```

<b>400 Bad Request</b>

If the `deck_id` specified is not a valid UUID, 400 is returned.

<b>404 Not Found</b>

If the `deck_id` specified does not exist, 404 is returned.

#### Example - reshuffle a deck

```
curl -X POST http://localhost:4000/api/v1alpha/decks/{deck_id}/shuffle
```

### Deleting a deck

```
DELETE /api/v1alpha/decks/:deck_id
```

#### Params

- `deck_id` (**required**) - the ID of the deck to delete.

#### Responses

<b>204 No Content</b>

The deck was deleted successfully. No response body is returned.

<b>400 Bad Request</b>

If the `deck_id` specified is not a valid UUID, 400 is returned.

<b>404 Not Found</b>

If the `deck_id` specified does not exist, 404 is returned.

#### Example - delete a deck

```
curl -X DELETE http://localhost:4000/api/v1alpha/decks/{deck_id}
```
