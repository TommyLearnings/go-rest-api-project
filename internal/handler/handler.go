package handler

import (
	"encoding/json"
	"github.com/TommyLearning/go-rest-api-project/internal/logger"
	"github.com/TommyLearning/go-rest-api-project/internal/store"
	"net/http"

	"github.com/google/uuid"
)

type NewsStorer interface {
	Create(store.News) (store.News, error)
	FindById(uuid.UUID) (store.News, error)
	FindAll() ([]store.News, error)
	UpdateById(store.News) error
	DeleteById(uuid.UUID) error
}

func PostNews(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.FromContext(r.Context())
		logger.Info("post news")

		var requestBody NewsPostReqBody
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			logger.Error("failed to decode request body", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		n, err := requestBody.Validate()
		if err != nil {
			logger.Error("failed to validate request body", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if _, err := ns.Create(n); err != nil {
			logger.Error("failed to create news", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func GetAllNews(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.FromContext(r.Context())
		logger.Info("get all news")
		news, err := ns.FindAll()
		if err != nil {
			logger.Error("failed to get all news", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		allNewsResponse := AllNewsResponse{News: news}
		if err := json.NewEncoder(w).Encode(allNewsResponse); err != nil {
			logger.Error("failed to encode response", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func GetNewsById(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.FromContext(r.Context())
		logger.Info("get news by id")
		newsID := r.PathValue("news_id")
		newsUUID, err := uuid.Parse(newsID)
		if err != nil {
			logger.Error("failed to parse news id", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		news, err := ns.FindById(newsUUID)
		if err != nil {
			logger.Error("failed to get news by id", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(news); err != nil {
			logger.Error("failed to encode response", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
}

func UpdateNewsById(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.FromContext(r.Context())
		logger.Info("update news by id")

		var newsReqBody NewsPostReqBody
		if err := json.NewDecoder(r.Body).Decode(&newsReqBody); err != nil {
			logger.Error("failed to decode the request", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		n, err := newsReqBody.Validate()
		if err != nil {
			logger.Error("failed to validate request body", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if err2 := ns.UpdateById(n); err2 != nil {
			logger.Error("failed to update news by id", "error", err2)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
}

func DeleteNewsById(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.FromContext(r.Context())
		newsID := r.PathValue("news_id")
		newsUUID, err := uuid.Parse(newsID)
		if err != nil {
			logger.Error("failed to parse news id", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := ns.DeleteById(newsUUID); err != nil {
			logger.Error("failed to delete news by id", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	}
}

type AllNewsResponse struct {
	News []store.News `json:"news"`
}
