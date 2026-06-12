package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type article struct {
	ID        string `json:"id"`
	Slug      string `json:"slug"`
	Thumbnail string `json:"thumbnail"`
}

type articleFile struct {
	Articles []article `json:"articles"`
}

func main() {
	jsonPath := flag.String("json", "mockdb/articles.json", "path to articles JSON with thumbnails")
	creds := flag.String("creds", "golang-blogs-firebase-adminsdk-fbsvc-64dacce61f.json", "Firebase service account JSON")
	project := flag.String("project", "golang-blogs", "Firebase project ID")
	collection := flag.String("collection", "articles", "Firestore collection name")
	flag.Parse()

	data, err := os.ReadFile(*jsonPath)
	if err != nil {
		log.Fatalf("read json: %v", err)
	}

	var file articleFile
	if err := json.Unmarshal(data, &file); err != nil {
		log.Fatalf("parse json: %v", err)
	}

	byID := make(map[string]string)
	for _, a := range file.Articles {
		if strings.TrimSpace(a.Thumbnail) == "" {
			continue
		}
		id := a.ID
		if id == "" {
			id = a.Slug
		}
		if id != "" {
			byID[id] = a.Thumbnail
		}
		if a.Slug != "" {
			byID[a.Slug] = a.Thumbnail
		}
	}

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, *project, option.WithCredentialsFile(*creds))
	if err != nil {
		log.Fatalf("firestore client: %v", err)
	}
	defer client.Close()

	col := client.Collection(*collection)
	iter := col.Documents(ctx)
	defer iter.Stop()

	patched := 0
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("read docs: %v", err)
		}

		thumb, ok := byID[doc.Ref.ID]
		if !ok {
			if slug, _ := doc.DataAt("slug"); slug != nil {
				if s, ok := slug.(string); ok {
					thumb, ok = byID[s]
				}
			}
		}
		if !ok || strings.TrimSpace(thumb) == "" {
			if cat, _ := doc.DataAt("category"); cat != nil {
				if c, ok := cat.(string); ok {
					thumb, ok = categoryThumbnail(c)
				}
			}
		}
		if strings.TrimSpace(thumb) == "" {
			thumb = "https://images.unsplash.com/photo-1621761199169-4d2c036feb09?auto=format&fit=crop&w=1200&q=80"
		}

		if _, err := doc.Ref.Set(ctx, map[string]any{"thumbnail": thumb}, firestore.MergeAll); err != nil {
			log.Fatalf("patch %s: %v", doc.Ref.ID, err)
		}
		fmt.Printf("patched %s\n", doc.Ref.ID)
		patched++
	}

	log.Printf("done: patched %d articles in %s/%s", patched, *project, *collection)
}

func categoryThumbnail(category string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(category)) {
	case "bitcoin":
		return "https://images.unsplash.com/photo-1518546305921-5a4ffdbb5feb?auto=format&fit=crop&w=1200&q=80", true
	case "ethereum":
		return "https://images.unsplash.com/photo-1639765487014-cfd9c42ad98d?auto=format&fit=crop&w=1200&q=80", true
	case "defi", "derivatives":
		return "https://images.unsplash.com/photo-1611974789855-9c2a0a7236a3?auto=format&fit=crop&w=1200&q=80", true
	default:
		return "", false
	}
}
