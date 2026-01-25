package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/TommyLearning/go-rest-api-project/internal/logger"
	"github.com/TommyLearning/go-rest-api-project/internal/news"

	"github.com/google/uuid"
)

//go:generate mockgen -source=handler.go -destination=mocks/handler.go -package=mockshandler

type NewsStorer interface {
	//Create(store.News) (store.News, error)
	//FindById(uuid.UUID) (store.News, error)
	//FindAll() ([]store.News, error)
	//UpdateById(store.News) error
	//DeleteById(uuid.UUID) error
	Create(context.Context, *news.Record) (*news.Record, error)
	FindById(context.Context, uuid.UUID) (*news.Record, error)
	FindAll(context.Context) ([]*news.Record, error)
	DeleteById(context.Context, uuid.UUID) error
	UpdateById(context.Context, uuid.UUID, *news.Record) error
}

func PostNews(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)
		log.Info("post news")

		var requestBody NewsPostReqBody
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			log.Error("failed to decode request body", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		n, err := requestBody.Validate()
		if err != nil {
			log.Error("failed to validate request body", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if _, err := ns.Create(ctx, n); err != nil {
			log.Error("failed to create news", "error", err)
			var dbErr *news.CustomError
			if errors.As(err, &dbErr) {
				w.WriteHeader(dbErr.HttpStatusCode())
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func GetAllNews(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)
		log.Info("get all news")
		n, err := ns.FindAll(ctx)
		if err != nil {
			log.Error("failed to get all news", "error", err)
			var dbErr *news.CustomError
			if errors.As(err, &dbErr) {
				w.WriteHeader(dbErr.HttpStatusCode())
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		allNewsResponse := AllNewsResponse{News: n}
		if err := json.NewEncoder(w).Encode(allNewsResponse); err != nil {
			log.Error("failed to encode response", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func GetNewsById(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		log := logger.FromContext(ctx)
		log.Info("get news by id")
		newsID := r.PathValue("news_id")
		newsUUID, err := uuid.Parse(newsID)
		if err != nil {
			log.Error("failed to parse news id", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		n, err := ns.FindById(ctx, newsUUID)
		if err != nil {
			log.Error("failed to get news by id", "error", err)
			var dbErr *news.CustomError
			if errors.As(err, &dbErr) {
				w.WriteHeader(dbErr.HttpStatusCode())
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(n); err != nil {
			log.Error("failed to encode response", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
}

func UpdateNewsById(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)
		log.Info("update news by id")

		var newsReqBody NewsPostReqBody
		if err := json.NewDecoder(r.Body).Decode(&newsReqBody); err != nil {
			log.Error("failed to decode the request", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		n, err := newsReqBody.Validate()
		if err != nil {
			log.Error("failed to validate request body", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if err2 := ns.UpdateById(ctx, n.Id, n); err2 != nil {
			log.Error("failed to update news by id", "error", err2)
			var dbErr *news.CustomError
			if errors.As(err, &dbErr) {
				w.WriteHeader(dbErr.HttpStatusCode())
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
}

func DeleteNewsById(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)
		newsID := r.PathValue("news_id")
		newsUUID, err := uuid.Parse(newsID)
		if err != nil {
			log.Error("failed to parse news id", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := ns.DeleteById(ctx, newsUUID); err != nil {
			log.Error("failed to delete news by id", "error", err)
			var dbErr *news.CustomError
			if errors.As(err, &dbErr) {
				w.WriteHeader(dbErr.HttpStatusCode())
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	}
}
