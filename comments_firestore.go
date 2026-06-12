package main

import (
	"context"
	"fmt"
	"sort"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type firestoreCommentStore struct {
	client     *firestore.Client
	collection string
}

func newFirestoreCommentStore(ctx context.Context, projectID, collection string) (*firestoreCommentStore, error) {
	opts, err := firestoreClientOptions()
	if err != nil {
		return nil, err
	}
	client, err := firestore.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, fmt.Errorf("firestore comment client: %w", err)
	}
	return &firestoreCommentStore{client: client, collection: collection}, nil
}

func (s *firestoreCommentStore) ListByArticle(ctx context.Context, slug string) ([]Comment, error) {
	iter := s.client.Collection(s.collection).
		Where("articleSlug", "==", slug).
		Documents(ctx)
	defer iter.Stop()

	comments := make([]Comment, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("firestore comments read: %w", err)
		}

		var c Comment
		if err := doc.DataTo(&c); err != nil {
			return nil, fmt.Errorf("firestore comment decode: %w", err)
		}
		if c.ID == "" {
			c.ID = doc.Ref.ID
		}
		comments = append(comments, c)
	}

	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt < comments[j].CreatedAt
	})
	return comments, nil
}

func (s *firestoreCommentStore) Create(ctx context.Context, c Comment) error {
	ref := s.client.Collection(s.collection).NewDoc()
	c.ID = ref.ID
	if c.CreatedAt == "" {
		c.CreatedAt = newComment(c.ArticleSlug, commentForm{
			Nickname: c.Nickname,
			Body:     c.Body,
		}).CreatedAt
	}
	if _, err := ref.Set(ctx, c); err != nil {
		return fmt.Errorf("firestore comment write: %w", err)
	}
	return nil
}

func (s *firestoreCommentStore) Close() error {
	return s.client.Close()
}
