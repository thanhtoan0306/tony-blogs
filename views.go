package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type ViewCounter interface {
	Record(ctx context.Context, pageKey string) (int64, error)
}

func newViewCounter(ctx context.Context) (ViewCounter, string, error) {
	if os.Getenv("DATA_SOURCE") == "json" {
		return newMemoryViewCounter(), "memory:page_views", nil
	}

	projectID := firestoreProjectID()
	collection := os.Getenv("PAGE_VIEWS_COLLECTION")
	if collection == "" {
		collection = "page_views"
	}
	counter, err := newFirestoreViewCounter(ctx, projectID, collection)
	if err != nil {
		return nil, "", err
	}
	return counter, fmt.Sprintf("firestore:%s/%s", projectID, collection), nil
}

func recordPageView(ctx context.Context, views ViewCounter, pageKey string) int64 {
	if views == nil {
		return 0
	}
	n, err := views.Record(ctx, pageKey)
	if err != nil {
		log.Printf("page view %s: %v", pageKey, err)
		return 0
	}
	return n
}

func todayVNDate() string {
	return time.Now().In(vnLocation()).Format("2006-01-02")
}

func viewDocID(date, pageKey string) string {
	return date + "_" + sanitizePageKey(pageKey)
}

func sanitizePageKey(key string) string {
	key = strings.ReplaceAll(key, "/", "_")
	key = strings.ReplaceAll(key, " ", "-")
	return key
}
