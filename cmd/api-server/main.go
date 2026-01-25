package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/TommyLearning/go-rest-api-project/internal/logger"
	"github.com/TommyLearning/go-rest-api-project/internal/news"
	"github.com/TommyLearning/go-rest-api-project/internal/postgres"
	"github.com/TommyLearning/go-rest-api-project/internal/router"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	db, err := postgres.NewDB(&postgres.Config{})

	if err != nil {
		log.Error("failed to connect to db", "error", err)
		os.Exit(1)
	}
	newsStore := news.NewStore(db)
	r := router.New(newsStore)

	wrapperRouter := logger.AddLoggerMid(log, logger.LoggerMid(r))

	log.Info("server starting on port 8080")

	if err := http.ListenAndServe(":8080", wrapperRouter); err != nil {
		log.Error("faild to start server", "error", err)
	}

}
