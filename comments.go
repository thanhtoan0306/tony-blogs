package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

type Comment struct {
	ID          string `json:"id" firestore:"id"`
	ArticleSlug string `json:"articleSlug" firestore:"articleSlug"`
	Nickname    string `json:"nickname" firestore:"nickname"`
	Body        string `json:"body" firestore:"body"`
	CreatedAt   string `json:"createdAt" firestore:"createdAt"`
}

type CommentStore interface {
	ListByArticle(ctx context.Context, slug string) ([]Comment, error)
	Create(ctx context.Context, c Comment) error
}

type commentForm struct {
	Nickname string
	Body     string
}

func newCommentStore(ctx context.Context) (CommentStore, string, error) {
	if os.Getenv("DATA_SOURCE") == "json" {
		return newMemoryCommentStore(), "memory:comments", nil
	}

	projectID := firestoreProjectID()
	collection := os.Getenv("COMMENTS_COLLECTION")
	if collection == "" {
		collection = "comments"
	}
	store, err := newFirestoreCommentStore(ctx, projectID, collection)
	if err != nil {
		return nil, "", err
	}
	return store, fmt.Sprintf("firestore:%s/%s", projectID, collection), nil
}

func validateCommentForm(form commentForm) (commentForm, error) {
	form.Nickname = strings.TrimSpace(form.Nickname)
	form.Body = strings.TrimSpace(form.Body)

	if len(form.Nickname) < 2 {
		return form, fmt.Errorf("nickname cần ít nhất 2 ký tự")
	}
	if len(form.Nickname) > 30 {
		return form, fmt.Errorf("nickname tối đa 30 ký tự")
	}
	if form.Body == "" {
		return form, fmt.Errorf("nội dung bình luận không được để trống")
	}
	if len(form.Body) > 1000 {
		return form, fmt.Errorf("nội dung tối đa 1000 ký tự")
	}
	return form, nil
}

func commentInitial(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "?"
	}
	r, _ := utf8.DecodeRuneInString(s)
	return strings.ToUpper(string(r))
}

func newComment(slug string, form commentForm) Comment {
	return Comment{
		ArticleSlug: slug,
		Nickname:    form.Nickname,
		Body:        form.Body,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
	}
}
