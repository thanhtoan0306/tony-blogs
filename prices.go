package main

import (
	"context"
	"sync"
	"time"
)

const pricePollInterval = time.Hour

type GoldQuote struct {
	Name       string
	BuyVND     int64
	SellVND    int64
	ChangeBuy  int64
	ChangeSell int64
	Unit       string
	Source     string
}

type FuelQuote struct {
	Name      string
	Region1   int
	Region2   int
	Unit      string
	Source    string
}

type OilQuote struct {
	Name      string
	Symbol    string
	PriceUSD  float64
	ChangeUSD float64
	Unit      string
	Source    string
}

type USDQuote struct {
	BuyVND  int64
	SellVND int64
	AvgVND  int64
	Unit    string
	Source  string
}

type BTCQuote struct {
	Name      string
	Symbol    string
	PriceUSD  float64
	PriceVND  int64
	ChangeUSD float64
	Unit      string
	Source    string
}

type VNPriceBoard struct {
	UpdatedAt  time.Time
	NextUpdate time.Time
	Gold       GoldQuote
	USD        USDQuote
	BTC        BTCQuote
	Gasoline   []FuelQuote
	Oil        []OilQuote
	Errors     []string
}

type priceBoard struct {
	mu   sync.RWMutex
	data VNPriceBoard
}

func newPriceBoard() *priceBoard {
	return &priceBoard{}
}

func (pb *priceBoard) Start(ctx context.Context) {
	pb.refresh(ctx)
	go func() {
		ticker := time.NewTicker(pricePollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				pb.refresh(ctx)
			}
		}
	}()
}

func (pb *priceBoard) Snapshot(ctx context.Context) VNPriceBoard {
	pb.mu.RLock()
	stale := pb.data.UpdatedAt.IsZero() || time.Since(pb.data.UpdatedAt) >= pricePollInterval
	pb.mu.RUnlock()
	if stale {
		pb.refresh(ctx)
	}

	pb.mu.RLock()
	defer pb.mu.RUnlock()
	return pb.data
}

func (pb *priceBoard) refresh(ctx context.Context) {
	board := VNPriceBoard{
		UpdatedAt:  time.Now().UTC(),
		NextUpdate: time.Now().UTC().Add(pricePollInterval),
	}

	if gold, err := fetchGoldVN(ctx); err != nil {
		board.Errors = append(board.Errors, "vàng: "+err.Error())
	} else {
		board.Gold = gold
	}

	if usd, err := fetchUSDVN(ctx); err != nil {
		board.Errors = append(board.Errors, "USD: "+err.Error())
	} else {
		board.USD = usd
	}

	if btc, err := fetchBTC(ctx); err != nil {
		board.Errors = append(board.Errors, "BTC: "+err.Error())
	} else {
		if board.USD.AvgVND > 0 {
			btc.PriceVND = int64(btc.PriceUSD*float64(board.USD.AvgVND) + 0.5)
		}
		board.BTC = btc
	}

	if gas, err := fetchGasVN(ctx); err != nil {
		board.Errors = append(board.Errors, "xăng dầu: "+err.Error())
	} else {
		board.Gasoline = gas
	}

	if oil, err := fetchOilWorld(ctx); err != nil {
		board.Errors = append(board.Errors, "dầu thế giới: "+err.Error())
	} else {
		board.Oil = oil
	}

	pb.mu.Lock()
	pb.data = board
	pb.mu.Unlock()
}
