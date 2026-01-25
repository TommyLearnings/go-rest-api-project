package router

import (
	"github.com/TommyLearning/go-rest-api-project/internal/handler"
	"github.com/TommyLearning/go-rest-api-project/internal/news"

	"net/http"
)

func New(ns *news.Store) *http.ServeMux {

	r := http.NewServeMux()

	r.HandleFunc("POST /news", handler.PostNews(ns))
	r.HandleFunc("GET /news", handler.GetAllNews(ns))
	r.HandleFunc("GET /news/{news_id}", handler.GetNewsById(ns))
	r.HandleFunc("PUT /news/{news_id}", handler.UpdateNewsById(ns))
	r.HandleFunc("DELETE /news/{news_id}", handler.DeleteNewsById(ns))

	return r
}
