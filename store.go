package main

import (
	"context"
	"fmt"
	"os"
)

// memoryStore implements ArticleManager for local JSON mode.

// ArticleStore loads and serves articles from Firestore by default (JSON via DATA_SOURCE=json).
type ArticleStore interface {
	All() []Article
	BySlug(slug string) (Article, bool)
}

func newArticleStore(ctx context.Context) (ArticleStore, string, error) {
	switch os.Getenv("DATA_SOURCE") {
	case "json":
		path := os.Getenv("ARTICLES_JSON")
		if path == "" {
			path = "mockdb/articles.json"
		}
		articles, err := loadArticlesFromJSON(path)
		if err != nil {
			return nil, "", err
		}
		return newMemoryStore(articles), path, nil
	default:
		projectID := firestoreProjectID()
		collection := os.Getenv("FIRESTORE_COLLECTION")
		if collection == "" {
			collection = "articles"
		}
		store, err := newFirestoreStore(ctx, projectID, collection)
		if err != nil {
			return nil, "", err
		}
		return store, fmt.Sprintf("firestore:%s/%s", projectID, collection), nil
	}
}

type memoryStore struct {
	articles []Article
	bySlug   map[string]Article
}

func newMemoryStore(articles []Article) *memoryStore {
	bySlug := make(map[string]Article, len(articles))
	for _, a := range articles {
		bySlug[a.Slug] = a
	}
	return &memoryStore{articles: articles, bySlug: bySlug}
}

func (s *memoryStore) All() []Article { return s.articles }

func (s *memoryStore) BySlug(slug string) (Article, bool) {
	a, ok := s.bySlug[slug]
	return a, ok
}

func (s *memoryStore) SetVisible(_ context.Context, slug string, visible bool) error {
	a, ok := s.bySlug[slug]
	if !ok {
		return fmt.Errorf("article %q not found", slug)
	}
	a.Visible = boolPtr(visible)
	for i := range s.articles {
		if s.articles[i].Slug == slug {
			s.articles[i] = a
			break
		}
	}
	s.bySlug[slug] = a
	return nil
}
