# golang-news

Server-rendered crypto news in Go. Articles are loaded from **Firestore** by default and rendered as HTML on every request.

## Quick start

```bash
go run .
```

Requires `golang-blogs-firebase-adminsdk-fbsvc-64dacce61f.json` in the project root (see Firebase below).

Open [http://127.0.0.1:8093](http://127.0.0.1:8093).

## Routes

| Path | Description |
|------|-------------|
| `GET /` | Home feed (today's stories, or latest if none today) |
| `GET /news/{slug}` | Single article |
| `GET /health` | JSON health check |
| `GET /static/global.css` | Styles |

## Firebase (Firestore)

Default data source. Place the service account JSON in the project root (never commit it — already in `.gitignore`).

Seed or refresh articles from local mock data:

```bash
go run ./cmd/seed-firestore
```

## Mock data (local dev)

Use JSON files instead of Firestore:

```bash
DATA_SOURCE=json go run .
```

`mockdb/articles.json` — HTML bodies (`bodyHtml`). Plain-text fallback: `data/articles.json` (`body` field).

```bash
DATA_SOURCE=json ARTICLES_JSON=data/articles.json go run .
```

## Env

| Variable | Default |
|----------|---------|
| `PORT` | `8093` |
| `DATA_SOURCE` | `firebase` (Firestore); set `json` for local files |
| `ARTICLES_JSON` | `mockdb/articles.json` (when `DATA_SOURCE=json`) |
| `FIREBASE_CREDENTIALS` | `golang-blogs-firebase-adminsdk-fbsvc-64dacce61f.json` |
| `FIREBASE_PROJECT_ID` | `golang-blogs` |
| `FIRESTORE_COLLECTION` | `articles` |
