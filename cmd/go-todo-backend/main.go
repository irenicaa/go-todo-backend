package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/irenicaa/go-todo-backend/gateways/db"
	"github.com/irenicaa/go-todo-backend/gateways/handlers"
	httputils "github.com/irenicaa/go-todo-backend/http-utils"
	usecases "github.com/irenicaa/go-todo-backend/use-cases"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	dbDSN, ok := os.LookupEnv("DB_DSN")
	if !ok {
		dbDSN = "postgresql://postgres:postgres@localhost:5432" +
			"/postgres?sslmode=disable"
	}
	flag.Parse()

	logger := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	storage, err := db.OpenDB(dbDSN)
	if err != nil {
		logger.Fatal(err)
	}

	handler := httputils.LoggingMiddleware(
		handlers.Router{
			BaseURL: "/api/v1",
			TodoRecord: handlers.TodoRecord{
				URLScheme: "http",
				UseCase: usecases.TodoRecord{
					Storage: storage,
				},
				Logger: logger,
			},
			Logger: logger,
		},
		logger,
		time.Now,
	)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		logger.Fatal(err)
	}
}
