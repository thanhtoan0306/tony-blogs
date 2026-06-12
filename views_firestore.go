package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

type firestoreViewCounter struct {
	client     *firestore.Client
	collection string
}

func newFirestoreViewCounter(ctx context.Context, projectID, collection string) (*firestoreViewCounter, error) {
	opts, err := firestoreClientOptions()
	if err != nil {
		return nil, err
	}
	client, err := firestore.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, fmt.Errorf("firestore view client: %w", err)
	}
	return &firestoreViewCounter{client: client, collection: collection}, nil
}

func (v *firestoreViewCounter) Record(ctx context.Context, pageKey string) (int64, error) {
	date := todayVNDate()
	ref := v.client.Collection(v.collection).Doc(viewDocID(date, pageKey))

	if _, err := ref.Set(ctx, map[string]any{
		"page":  pageKey,
		"date":  date,
		"count": firestore.Increment(1),
	}, firestore.MergeAll); err != nil {
		return 0, fmt.Errorf("firestore view increment: %w", err)
	}

	snap, err := ref.Get(ctx)
	if err != nil {
		return 0, fmt.Errorf("firestore view read: %w", err)
	}
	count, err := snap.DataAt("count")
	if err != nil {
		return 0, fmt.Errorf("firestore view count field: %w", err)
	}
	return toInt64(count), nil
}

func (v *firestoreViewCounter) Close() error {
	return v.client.Close()
}

func toInt64(v any) int64 {
	switch n := v.(type) {
	case int64:
		return n
	case int:
		return int64(n)
	case float64:
		return int64(n)
	default:
		return 0
	}
}
