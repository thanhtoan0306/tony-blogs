package main

import (
	"context"
	"sort"
	"sync"
)

type memoryCommentStore struct {
	mu       sync.Mutex
	comments []Comment
	nextID   int
}

func newMemoryCommentStore() *memoryCommentStore {
	return &memoryCommentStore{nextID: 1}
}

func (s *memoryCommentStore) ListByArticle(_ context.Context, slug string) ([]Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]Comment, 0)
	for _, c := range s.comments {
		if c.ArticleSlug == slug {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt < out[j].CreatedAt
	})
	return out, nil
}

func (s *memoryCommentStore) Create(_ context.Context, c Comment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	c.ID = itoa(s.nextID)
	s.nextID++
	s.comments = append(s.comments, c)
	return nil
}
