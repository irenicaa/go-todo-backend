package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/irenicaa/go-http-utils/middlewares"
	"github.com/irenicaa/go-todo-backend/v2/gateways/db"
	"github.com/irenicaa/go-todo-backend/v2/gateways/handlers"
	usecases "github.com/irenicaa/go-todo-backend/v2/use-cases"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	dbDSN, ok := os.LookupEnv("DB_DSN")
	if !ok {
		dbDSN = db.DefaultDataSourceName
	}
	flag.Parse()

	logger := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	dbPool, err := db.OpenDB(dbDSN)
	if err != nil {
		logger.Fatal(err)
	}

	handler := middlewares.LoggingMiddleware(
		middlewares.CORSMiddleware(handlers.Router{
			BaseURL: "/api/v1",
			TodoRecord: handlers.TodoRecord{
				URLScheme: "http",
				UseCase: usecases.TodoRecord{
					Storage: db.NewTodoRecord(dbPool),
				},
				Logger: logger,
			},
			Logger: logger,
		}),
		logger,
		time.Now,
	)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		logger.Fatal(err)
	}
}
