package main

import (
	"context"
	"fmt"
	"sort"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type firestoreStore struct {
	client     *firestore.Client
	collection string
	articles   []Article
	bySlug     map[string]Article
}

func newFirestoreStore(ctx context.Context, credsPath, projectID, collection string) (*firestoreStore, error) {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(credsPath))
	if err != nil {
		return nil, fmt.Errorf("firestore client: %w", err)
	}

	store := &firestoreStore{
		client:     client,
		collection: collection,
	}
	if err := store.reload(ctx); err != nil {
		_ = client.Close()
		return nil, err
	}
	return store, nil
}

func (s *firestoreStore) reload(ctx context.Context) error {
	iter := s.client.Collection(s.collection).Documents(ctx)
	defer iter.Stop()

	articles := make([]Article, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("firestore read: %w", err)
		}

		var a Article
		if err := doc.DataTo(&a); err != nil {
			return fmt.Errorf("firestore decode %s: %w", doc.Ref.ID, err)
		}
		if a.ID == "" {
			a.ID = doc.Ref.ID
		}
		if a.Slug == "" {
			a.Slug = doc.Ref.ID
		}
		articles = append(articles, a)
	}

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].PublishedAt > articles[j].PublishedAt
	})

	bySlug := make(map[string]Article, len(articles))
	for _, a := range articles {
		bySlug[a.Slug] = a
	}
	s.articles = articles
	s.bySlug = bySlug
	return nil
}

func (s *firestoreStore) All() []Article { return s.articles }

func (s *firestoreStore) BySlug(slug string) (Article, bool) {
	a, ok := s.bySlug[slug]
	return a, ok
}

func (s *firestoreStore) Create(ctx context.Context, a Article) error {
	id := a.ID
	if id == "" {
		id = a.Slug
	}
	if id == "" {
		return fmt.Errorf("article id is required")
	}
	a.ID = id
	if a.Slug == "" {
		a.Slug = id
	}

	if _, err := s.client.Collection(s.collection).Doc(id).Set(ctx, a); err != nil {
		return fmt.Errorf("firestore write: %w", err)
	}
	return s.reload(ctx)
}

func (s *firestoreStore) Close() error {
	return s.client.Close()
}
