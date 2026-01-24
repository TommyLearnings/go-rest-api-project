package handler_test

import (
	"errors"
	"github.com/TommyLearning/go-rest-api-project/internal/handler"
	"github.com/TommyLearning/go-rest-api-project/internal/store"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func Test_PostNews(t *testing.T) {
	testCases := []struct {
		name           string
		body           io.Reader
		store          handler.NewsStorer
		expectedStatus int
	}{
		{
			name:           "invalid request body json",
			body:           strings.NewReader(`{`),
			store:          mockNewsStore{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid request body",
			body: strings.NewReader(`
				{
				"id" : "3b082d9d-1dc7-4d1f-907e-50d449a03d45",
				"author": "code learn",
				"title": "first news",
				"summary": "first news post",
				"created_at": "2024-04-07T05:13:27+00:00",
				"source": "https://example.com"
				}`),
			store:          mockNewsStore{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "db error",
			body: strings.NewReader(`
				{
				"id" : "3b082d9d-1dc7-4d1f-907e-50d449a03d45",
				"author": "code learn",
				"title": "first news",
				"summary": "first news post",
				"created_at": "2024-04-07T05:13:27+00:00",
				"source": "https://example.com",
				"tags": ["politics"]
				}`),
			store:          mockNewsStore{errState: true},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "success",
			body: strings.NewReader(`
				{
				"id" : "3b082d9d-1dc7-4d1f-907e-50d449a03d45",
				"author": "code learn",
				"title": "first news",
				"summary": "first news post",
				"created_at": "2024-04-07T05:13:27+00:00",
				"source": "https://example.com",
				"tags": ["politics"]
				}`),
			store:          mockNewsStore{},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//arrange
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", tc.body)

			//act
			handler.PostNews(tc.store)(w, r)

			//Assert
			if w.Result().StatusCode != tc.expectedStatus {
				t.Errorf("expected status %d got %d", tc.expectedStatus, w.Result().StatusCode)
			}
		})
	}
}

func Test_GetAllNews(t *testing.T) {
	testCases := []struct {
		name           string
		store          handler.NewsStorer
		expectedStatus int
	}{{
		name:           "db error",
		store:          mockNewsStore{errState: true},
		expectedStatus: http.StatusInternalServerError,
	},
		{
			name:           "success",
			store:          mockNewsStore{},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//arrange
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", nil)

			//act
			handler.GetAllNews(tc.store)(w, r)

			//Assert
			if w.Result().StatusCode != tc.expectedStatus {
				t.Errorf("expected status %d got %d", tc.expectedStatus, w.Result().StatusCode)
			}
		})
	}
}

func TestGetNewsById(t *testing.T) {
	testCases := []struct {
		name           string
		store          handler.NewsStorer
		newsId         string
		expectedStatus int
	}{
		{
			name:           "invalid news id",
			store:          mockNewsStore{},
			newsId:         "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "db error",
			store:          mockNewsStore{errState: true},
			newsId:         uuid.NewString(),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "success",
			store:          mockNewsStore{},
			newsId:         uuid.NewString(),
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//arrange
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", nil)
			r.SetPathValue("news_id", tc.newsId)
			//act
			handler.GetNewsById(tc.store)(w, r)

			//Assert
			if w.Result().StatusCode != tc.expectedStatus {
				t.Errorf("expected status %d got %d", tc.expectedStatus, w.Result().StatusCode)
			}
		})
	}
}

func TestUpdateNewsById(t *testing.T) {
	testCases := []struct {
		name           string
		body           io.Reader
		store          handler.NewsStorer
		expectedStatus int
	}{
		{
			name:           "invalid request body json",
			body:           strings.NewReader(`{`),
			store:          mockNewsStore{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid request body",
			body: strings.NewReader(`
				{
				"id" : "3b082d9d-1dc7-4d1f-907e-50d449a03d45",
				"author": "code learn",
				"title": "first news",
				"summary": "first news post",
				"created_at": "2024-04-07T05:13:27+00:00",
				"source": "https://example.com"
				}`),
			store:          mockNewsStore{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "db error",

			body: strings.NewReader(`
			{ 
			"id" : "3b082d9d-1dc7-4d1f-907e-50d449a03d45", 
			"author": "code learn", 
			"content": "news content",
			"title": "first news", 
			"summary": "first news post", 
			"created_at": "2024-04-07T05:13:27+00:00", 
			"source": "https://example.com",
			"tags": ["politics"]
			}`),
			store:          mockNewsStore{errState: true},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "success",
			body: strings.NewReader(`
				{
				"id" : "3b082d9d-1dc7-4d1f-907e-50d449a03d45",
				"author": "code learn",
				"title": "first news",
				"summary": "first news post",
				"created_at": "2024-04-07T05:13:27+00:00",
				"source": "https://example.com",
				"tags": ["politics"]
				}
			)`),
			store:          mockNewsStore{},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//arrange
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", tc.body)

			//act
			handler.UpdateNewsById(tc.store)(w, r)

			//Assert
			if w.Result().StatusCode != tc.expectedStatus {
				t.Errorf("expected status %d got %d", tc.expectedStatus, w.Result().StatusCode)
			}
		})
	}
}

func TestDeleteNewsById(t *testing.T) {
	testCases := []struct {
		name           string
		store          handler.NewsStorer
		newsId         string
		expectedStatus int
	}{
		{
			name:           "invalid news id",
			store:          mockNewsStore{},
			newsId:         "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "db error",
			store:          mockNewsStore{errState: true},
			newsId:         uuid.NewString(),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "success",
			store:          mockNewsStore{},
			newsId:         uuid.NewString(),
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//arrange
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", nil)
			r.SetPathValue("news_id", tc.newsId)
			//act
			handler.DeleteNewsById(tc.store)(w, r)

			//Assert
			if w.Result().StatusCode != tc.expectedStatus {
				t.Errorf("expected status %d got %d", tc.expectedStatus, w.Result().StatusCode)
			}
		})
	}
}

type mockNewsStore struct {
	errState bool
}

func (m mockNewsStore) Create(_ store.News) (news store.News, err error) {
	if m.errState {
		return news, errors.New("error")
	}
	return news, nil
}

func (m mockNewsStore) FindById(_ uuid.UUID) (news store.News, err error) {
	if m.errState {
		return news, errors.New("error")
	}
	return news, nil
}

func (m mockNewsStore) FindAll() (news []store.News, err error) {
	if m.errState {
		return news, errors.New("error")
	}
	return news, nil
}

func (m mockNewsStore) UpdateById(_ store.News) error {
	if m.errState {
		return errors.New("error")
	}
	return nil
}

func (m mockNewsStore) DeleteById(_ uuid.UUID) error {
	if m.errState {
		return errors.New("error")
	}
	return nil
}
