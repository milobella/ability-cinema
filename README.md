# Cinema ability
Milobella ability to know about series & movies.

## Features
- See what is coming out or soon in your favorite theater;
- Manage your favorite theater.

## Prerequisites

- Having ``golang`` installed [instructions](https://golang.org/doc/install)
- Having ``go dep`` installed [instructions](https://golang.github.io/dep/docs/installation.html)

## Build

```bash
$ go build -o bin/ability cmd/ability/main.go
```

## Run

```bash
$ bin/ability
```

## Requests example

#### Get the movies for this evening in favorite theater.

```bash
$ curl -i -X POST http://localhost:4444/resolve -d '{"nlu":{"BestIntent": "LAST_SHOWTIME"}}'
HTTP/1.1 200 OK
Date: Sat, 01 Jun 2019 09:17:16 GMT
Content-Type: text/plain; charset=utf-8
Transfer-Encoding: chunked

{
	"nlg": {
		"sentence": "Here are the movies in {{theater}} this evening, in the {{location}}'s theater",
		"params": [{
			"name": "theater",
			"value": "La Strada",
			"type": "string"
		}, {
			"name": "location",
			"value": "Mouans-Sartoux",
			"type": "string"
		}]
	},
	"visu": [{
		"title": "\"Toy Story 4\"",
		"display": "\"Séances du dimanche 23 juin 2019 : 11:00 (film à 11:10)\""
	}, ... ]
}
```

## CHANGELOGS
- [Application changelog](./CHANGELOG.md)
- [Helm chart changelog](./helm/ability-cinema/CHANGELOG.md)
