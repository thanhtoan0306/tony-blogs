package main

import (
	"encoding/json"
	"net/http"
)

type pricesPageData struct {
	Title  string
	Board  VNPriceBoard
	PollMs int
}

func handlePricesPage(board *priceBoard) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=300")
		render(w, "prices.html", pricesPageData{
			Title:  "Giá Việt Nam · Crypto Today",
			Board:  board.Snapshot(r.Context()),
			PollMs: int(pricePollInterval.Milliseconds()),
		})
	}
}

func handlePricesAPI(board *priceBoard) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=300")
		_ = json.NewEncoder(w).Encode(board.Snapshot(r.Context()))
	}
}
