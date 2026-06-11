package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

type article struct {
	ID          string   `json:"id"`
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Summary     string   `json:"summary"`
	Body        string   `json:"body"`
	BodyHTML    string   `json:"bodyHtml"`
	Category    string   `json:"category"`
	Author      string   `json:"author"`
	PublishedAt string   `json:"publishedAt"`
	Tags        []string `json:"tags"`
}

type articleFile struct {
	Articles []article `json:"articles"`
}

func main() {
	jsonPath := flag.String("json", "mockdb/articles.json", "path to articles JSON")
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

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, *project, option.WithCredentialsFile(*creds))
	if err != nil {
		log.Fatalf("firestore client: %v", err)
	}
	defer client.Close()

	col := client.Collection(*collection)
	for _, a := range file.Articles {
		id := a.ID
		if id == "" {
			id = a.Slug
		}
		if id == "" {
			log.Fatal("article missing id and slug")
		}
		if _, err := col.Doc(id).Set(ctx, a); err != nil {
			log.Fatalf("write %s: %v", id, err)
		}
		fmt.Printf("seeded %s\n", id)
	}

	log.Printf("done: %d articles -> %s/%s", len(file.Articles), *project, *collection)
}
