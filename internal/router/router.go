package router

import (
	"net/http"

	"github.com/TommyLearning/go-rest-api-project/internal/handler"
)

func New(ns handler.NewsStorer) *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("POST /news", handler.PostNews(ns))
	r.HandleFunc("GET /news", handler.GetAllNews(ns))
	r.HandleFunc("GET /news/{news_id}", handler.GetNewsById(ns))
	r.HandleFunc("PUT /news/{news_id}", handler.UpdateNewsById(ns))
	r.HandleFunc("DELETE /news/{news_id}", handler.DeleteNewsById(ns))

	return r
}
