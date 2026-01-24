package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/TommyLearning/go-rest-api-project/internal/logger"
	"github.com/TommyLearning/go-rest-api-project/internal/router"
	"github.com/TommyLearning/go-rest-api-project/internal/store"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	r := router.New(store.New())
	wrapperRouter := logger.AddLoggerMid(log, logger.LoggerMid(r))

	log.Info("server starting on port 8080")

	if err := http.ListenAndServe(":8080", wrapperRouter); err != nil {
		log.Error("faild to start server", "error", err)
	}

}
