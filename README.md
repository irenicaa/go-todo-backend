# go-todo-backend

[![GoDoc](https://godoc.org/github.com/irenicaa/go-todo-backend?status.svg)](https://godoc.org/github.com/irenicaa/go-todo-backend)
[![Go Report Card](https://goreportcard.com/badge/github.com/irenicaa/go-todo-backend)](https://goreportcard.com/report/github.com/irenicaa/go-todo-backend)
[![Build Status](https://travis-ci.com/irenicaa/go-todo-backend.svg?branch=master)](https://travis-ci.com/irenicaa/go-todo-backend)
[![codecov](https://codecov.io/gh/irenicaa/go-todo-backend/branch/master/graph/badge.svg)](https://codecov.io/gh/irenicaa/go-todo-backend)

The web service that implements specs of the [Todo-Backend](https://www.todobackend.com/) project with some improvements.

## Installation

```
$ go get github.com/irenicaa/go-todo-backend/...
```

## Migrations

```
$ docker run -v $(pwd)/migrations:/migrations --network host migrate/migrate:v4.14.1 \
  -verbose -path=/migrations -database postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable up
```

## Usage

```
$ go-todo-backend -h | -help | --help
```

Options:

- `-h`, `-help`, `--help` &mdash; show the help message and exit.

Environment variables:

- `PORT` &mdash; server port (default: `8080`);
- `DB_DSN` &mdash; DB connection string (default: `postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable`).

## Testing

Running of the unit tests:

```
$ go test -race -cover ./...
```

Running of both the unit and integration tests:

```
$ go test -race -cover -tags integration ./...
```

Running the tests of the [Todo-Backend](https://www.todobackend.com/) project: follow link https://www.todobackend.com/specs/index.html?http://localhost:8080/api/v1/todos

## Docs

- [swagger.yaml](docs/swagger.yaml) &mdash; Swagger definition of the server API
- [postman_collection.json](docs/postman_collection.json) &mdash; Postman collection of the server API

## Output Example

```
2021/04/16 20:00:14.976714 POST /api/v1/todos 22.803305ms
2021/04/16 20:00:14.978912 GET /api/v1/todos/1018 1.067026ms
2021/04/16 20:00:14.998649 POST /api/v1/todos 17.712474ms
2021/04/16 20:00:15.009770 PUT /api/v1/todos/1019 9.955194ms
2021/04/16 20:00:15.014229 GET /api/v1/todos/1019 729.329µs
2021/04/16 20:00:15.020878 POST /api/v1/todos 5.521073ms
2021/04/16 20:00:15.031938 PATCH /api/v1/todos/1020 9.708337ms
2021/04/16 20:00:15.037017 GET /api/v1/todos/1020 869.903µs
2021/04/16 20:00:15.042976 DELETE /api/v1/todos 3.620271ms
2021/04/16 20:00:15.053973 POST /api/v1/todos 10.003377ms
2021/04/16 20:00:15.065116 POST /api/v1/todos 9.891636ms
2021/04/16 20:00:15.076273 POST /api/v1/todos 10.366496ms
2021/04/16 20:00:15.087464 POST /api/v1/todos 9.810042ms
2021/04/16 20:00:15.098607 POST /api/v1/todos 9.032801ms
2021/04/16 20:00:15.109738 POST /api/v1/todos 9.984774ms
2021/04/16 20:00:15.120763 POST /api/v1/todos 9.961407ms
2021/04/16 20:00:15.131872 POST /api/v1/todos 8.161943ms
2021/04/16 20:00:15.143359 POST /api/v1/todos 10.464461ms
2021/04/16 20:00:15.154795 POST /api/v1/todos 10.100347ms
2021/04/16 20:00:15.165769 POST /api/v1/todos 8.72751ms
2021/04/16 20:00:15.167434 GET /api/v1/todos 393.922µs
2021/04/16 20:00:15.348377 GET /api/v1/todos?minimal_date=2006-01-07 1.575302ms
2021/04/16 20:00:15.428437 GET /api/v1/todos?maximal_date=2006-01-07 1.279132ms
2021/04/16 20:00:15.508703 GET /api/v1/todos?maximal_date=2006-01-09&minimal_date=2006-01-05 1.501993ms
2021/04/16 20:00:15.588316 GET /api/v1/todos?title_fragment=even 1.352159ms
2021/04/16 20:00:15.669183 GET /api/v1/todos?page=3&page_size=2 1.303843ms
2021/04/16 20:00:15.814140 GET /api/v1/todos?maximal_date=2006-01-19&minimal_date=2006-01-05&page=3&page_size=2&title_fragment=even 1.50664ms
2021/04/16 20:00:15.819337 POST /api/v1/todos 7.238819ms
2021/04/16 20:00:15.822753 DELETE /api/v1/todos/1032 9.564257ms
2021/04/16 20:00:15.825053 unable to get the to-do record: sql: no rows in result set
2021/04/16 20:00:15.825092 GET /api/v1/todos/1032 765.518µs
```

## License

The MIT License (MIT)

Copyright &copy; 2021 irenica
