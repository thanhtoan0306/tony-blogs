package main

import (
	"encoding/json"
	"fmt"
	"os"

	"google.golang.org/api/option"
)

func firestoreClientOptions() ([]option.ClientOption, error) {
	if raw := os.Getenv("FIREBASE_CREDENTIALS_JSON"); raw != "" {
		if !json.Valid([]byte(raw)) {
			return nil, fmt.Errorf("FIREBASE_CREDENTIALS_JSON is not valid JSON")
		}
		return []option.ClientOption{option.WithCredentialsJSON([]byte(raw))}, nil
	}

	if path := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); path != "" {
		if len(path) > 0 && path[0] == '{' {
			if !json.Valid([]byte(path)) {
				return nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS JSON is invalid")
			}
			return []option.ClientOption{option.WithCredentialsJSON([]byte(path))}, nil
		}
		if _, err := os.Stat(path); err == nil {
			return []option.ClientOption{option.WithCredentialsFile(path)}, nil
		}
	}

	credsPath := os.Getenv("FIREBASE_CREDENTIALS")
	if credsPath == "" {
		credsPath = "golang-blogs-firebase-adminsdk-fbsvc-64dacce61f.json"
	}
	if _, err := os.Stat(credsPath); err != nil {
		return nil, fmt.Errorf(
			"firestore credentials not found: set FIREBASE_CREDENTIALS_JSON on Vercel or place %s locally",
			credsPath,
		)
	}
	return []option.ClientOption{option.WithCredentialsFile(credsPath)}, nil
}

func firestoreProjectID() string {
	if id := os.Getenv("FIREBASE_PROJECT_ID"); id != "" {
		return id
	}
	if raw := os.Getenv("FIREBASE_CREDENTIALS_JSON"); raw != "" {
		if id := projectIDFromJSON(raw); id != "" {
			return id
		}
	}
	return "golang-blogs"
}

func projectIDFromJSON(raw string) string {
	var creds struct {
		ProjectID string `json:"project_id"`
	}
	if err := json.Unmarshal([]byte(raw), &creds); err != nil {
		return ""
	}
	return creds.ProjectID
}
