# golang-news

Server-rendered crypto news in Go. Articles are loaded from a mock JSON file and rendered as HTML on every request.

## Quick start

```bash
cd june/golang-news
go run .
```

Open [http://127.0.0.1:8093](http://127.0.0.1:8093).

## Routes

| Path | Description |
|------|-------------|
| `GET /` | Home feed (today's stories, or latest if none today) |
| `GET /news/{slug}` | Single article |
| `GET /health` | JSON health check |
| `GET /static/global.css` | Styles |

## Mock data

`mockdb/articles.json` — article bodies stored as HTML (`bodyHtml` field).

Plain-text fallback: `data/articles.json` (`body` field, split into paragraphs).

```bash
ARTICLES_JSON=data/articles.json go run .
```

## Env

| Variable | Default |
|----------|---------|
| `PORT` | `8093` |
| `ARTICLES_JSON` | `mockdb/articles.json` |
