Card Game API
==============

A simple RESTful API for creating, shuffling, and drawing cards from a deck of playing cards. Built with Go and SQLite, this project provides an easy-to-use interface for managing decks of cards for various applications such as card games, simulations, or statistical analysis.

### Features
* Create a standard 52-card deck or a custom deck with specific cards
* Shuffle the deck with a random order of cards
* Draw any number of cards from the deck
* Store decks in an SQLite database for easy management and retrieval
* Lightweight and efficient implementation using Go

### Getting Started
#### Prerequisites
* Go 1.16 or later
* SQLite 3
#### Setting Up the Project
1. Clone the repository:
```bash
git clone https://github.com/natemago/card-games-api.git
```
2. Change to the project directory:
```bash
cd card-api
```
3. Install dependencies:
```bash
go mod tidy
```
4. Run the project
```bash
go run main.go
```
5. To run tests
```bash
go test ./tests
```

The API server should now be running on `http://localhost:8080`

## API Documentation

### 1. Create a Deck
Endpoint: `/deck`

Method: `POST`

**Query Parameters:**

> `shuffled`: (optional) true to return a shuffled deck, false or omitted for an unshuffled deck.
> `cards`: (optional) A comma-separated list of card codes to create a custom deck. Example: AS,KH,2D,JC,10C

**Success Response:**
Code: `200 OK`
Content:  _A JSON object containing the deck ID, remaining card count, shuffled status, and an array of cards._
Example: `/deck?shuffled=true&cards=AS,KH,2D`
```json
{
	"deck_id": "33636e44-41ce-4383-baff-70615eb7339f",
	"shuffled": true,
	"remaining": 3,
	"cards": [
		{
			"value": "2",
			"suit": "DIAMONDS",
			"code": "2D"
		},
		{
			"value": "KING",
			"suit": "HEARTS",
			"code": "KH"
		},
		{
			"value": "ACE",
			"suit": "SPADES",
			"code": "AS"
		}
	]
}
```

**Error Response:**

Code: `400 BAD REQUEST`
Content: _A JSON object with an error message indicating the issue with the request._
Example: `/deck?shuffled=true&cards=AS,KH,2D,jj,Kh`
```json
{
	"message": "invalid cards values: JJ, KH (duplicate)"
}
```

### 2. Draw Cards
Endpoint: `/deck/:deck_id/draw`

Method: `POST`

**URL Parameters:**

> `deck_id`: The ID of the deck to draw cards from.

**Query Parameters:**

> `count`: The number of cards to draw from the deck.


**Success Response:**
Code: `200 OK`
Content: _A JSON object containing the deck ID, remaining card count, and an array of drawn cards_.
Example: `deck/1a51c5b6-ec0c-4d5f-9f64-e6bae4d57780/draw?count=2`
```json
{
	"cards": [
		{
			"value": "5",
			"suit": "DIAMONDS",
			"code": "5D"
		},
		{
			"value": "4",
			"suit": "DIAMONDS",
			"code": "4D"
		}
	]
}
```

**Error Responses:**

> Code: `400 BAD REQUEST`
> Content: _A JSON object with an error message indicating the issue with the request._
> Code: `404 NOT FOUND`
> Content: _A JSON object with an error message indicating that the deck was not found._


### 3. Get Deck
Endpoint: `/deck/:deck_id`

Method: `GET`

**URL Parameters:**

> `deck_id`: The ID of the deck to retrieve.

**Success Response:**
Code: `200 OK`
Content: `A JSON object containing the deck ID, remaining card count, shuffled status, and an array of cards.`
Example: `/deck/336db108-2b9b-474f-98b0-3c8537fa2eb4`
```json
{
	{
	"deck_id": "336db108-2b9b-474f-98b0-3c8537fa2eb4",
	"shuffled": true,
	"remaining": 2,
	"cards": [
		{
			"value": "KING",
			"suit": "HEARTS",
			"code": "KH"
		},
		{
			"value": "5",
			"suit": "DIAMONDS",
			"code": "5D"
		}
	]
}
}
```

**Error Response:**
> Code: `404 NOT FOUND`
> Content: _A JSON object with an error message indicating that the deck was not found._

