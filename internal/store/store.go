package store

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type Store struct {
	l    sync.Mutex
	news []News
}

func New() *Store {
	return &Store{
		l:    sync.Mutex{},
		news: []News{},
	}
}

func (s *Store) Create(news News) (News, error) {

	s.l.Lock()
	defer s.l.Unlock()
	news.ID = uuid.New()
	s.news = append(s.news, news)
	return news, nil
}

func (s *Store) FindAll() ([]News, error) {
	s.l.Lock()
	defer s.l.Unlock()
	return s.news, nil
}

func (s *Store) FindById(id uuid.UUID) (News, error) {
	s.l.Lock()
	defer s.l.Unlock()
	for _, news := range s.news {
		if news.ID == id {
			return news, nil
		}
	}
	return News{}, errors.New("news not found")
}

func (s *Store) DeleteById(id uuid.UUID) error {
	s.l.Lock()
	defer s.l.Unlock()
	idx := func(id uuid.UUID) int {
		for i, n := range s.news {
			if n.ID == id {
				return i
			}
		}
		return -1
	}(id)
	if idx == -1 {
		return errors.New("news not found")
	}
	s.news = append(s.news[:idx], s.news[idx+1:]...)
	return nil
}

func (s *Store) UpdateById(news News) error {
	s.l.Lock()
	defer s.l.Unlock()
	for idx, n := range s.news {
		if n.ID == news.ID {
			s.news[idx] = news
			return nil
		}
	}
	return errors.New("news not found")
}
