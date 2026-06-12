package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var webgiaFuelRow = regexp.MustCompile(
	`<tr><th><a[^>]*>([^<]+)</a></th><td class="text-right">([\d.]+)</td><td class="text-right">([\d.]+)</td></tr>`,
)

var (
	webgiaUSDAvgRate = regexp.MustCompile(`data-mo="VND"\s+data-rate="(\d+)"`)
	webgiaUSDBuy     = regexp.MustCompile(`id="wdg-usd-m">([^<]+)<`)
	webgiaUSDSell    = regexp.MustCompile(`id="wdg-usd-b">([^<]+)<`)
)

func fetchGoldVN(ctx context.Context) (GoldQuote, error) {
	body, err := httpGet(ctx, "https://www.vang.today/api/prices?type=SJL1L10")
	if err != nil {
		return GoldQuote{}, err
	}

	var resp struct {
		Success    bool   `json:"success"`
		Name       string `json:"name"`
		Buy        int64  `json:"buy"`
		Sell       int64  `json:"sell"`
		ChangeBuy  int64  `json:"change_buy"`
		ChangeSell int64  `json:"change_sell"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return GoldQuote{}, fmt.Errorf("parse gold json: %w", err)
	}
	if !resp.Success {
		return GoldQuote{}, fmt.Errorf("gold api returned unsuccessful response")
	}

	name := resp.Name
	if name == "" {
		name = "SJC 9999"
	}

	return GoldQuote{
		Name:       name,
		BuyVND:     resp.Buy,
		SellVND:    resp.Sell,
		ChangeBuy:  resp.ChangeBuy,
		ChangeSell: resp.ChangeSell,
		Unit:       "₫/lượng",
		Source:     "vang.today",
	}, nil
}

func fetchUSDVN(ctx context.Context) (USDQuote, error) {
	body, err := httpGet(ctx, "https://webgia.com/ngoai-te/usd/")
	if err != nil {
		return USDQuote{}, err
	}
	html := string(body)

	avgMatch := webgiaUSDAvgRate.FindStringSubmatch(html)
	if len(avgMatch) < 2 {
		return USDQuote{}, fmt.Errorf("usd average rate not found on webgia")
	}
	avg, err := strconv.ParseInt(avgMatch[1], 10, 64)
	if err != nil {
		return USDQuote{}, fmt.Errorf("parse usd average: %w", err)
	}

	buy := avg
	if m := webgiaUSDBuy.FindStringSubmatch(html); len(m) >= 2 {
		if v, err := parseWebgiaForex(m[1]); err == nil && v > 0 {
			buy = v
		}
	}

	sell := avg
	if m := webgiaUSDSell.FindStringSubmatch(html); len(m) >= 2 {
		if v, err := parseWebgiaForex(m[1]); err == nil && v > 0 {
			sell = v
		}
	}

	return USDQuote{
		BuyVND:  buy,
		SellVND: sell,
		AvgVND:  avg,
		Unit:    "₫/USD",
		Source:  "webgia.com",
	}, nil
}

func parseWebgiaForex(raw string) (int64, error) {
	raw = strings.TrimSpace(raw)
	raw = strings.ReplaceAll(raw, ".", "")
	raw = strings.ReplaceAll(raw, ",", ".")
	f, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, err
	}
	return int64(f + 0.5), nil
}

func fetchGasVN(ctx context.Context) ([]FuelQuote, error) {
	body, err := httpGet(ctx, "https://webgia.com/gia-xang-dau/petrolimex/")
	if err != nil {
		return nil, err
	}

	matches := webgiaFuelRow.FindAllStringSubmatch(string(body), -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("no petrolimex rows found on webgia")
	}

	out := make([]FuelQuote, 0, len(matches))
	for _, m := range matches {
		name := strings.TrimSpace(m[1])
		out = append(out, FuelQuote{
			Name:    name,
			Region1: parseWebgiaPrice(m[2]),
			Region2: parseWebgiaPrice(m[3]),
			Unit:    "₫/lít",
			Source:  "webgia.com / Petrolimex",
		})
	}
	return out, nil
}

func fetchOilWorld(ctx context.Context) ([]OilQuote, error) {
	symbols := []struct {
		Symbol string
		Name   string
	}{
		{"BZ=F", "Brent crude"},
		{"CL=F", "WTI crude"},
	}

	out := make([]OilQuote, 0, len(symbols))
	var errs []string
	for _, s := range symbols {
		q, err := fetchYahooOil(ctx, s.Symbol, s.Name)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		out = append(out, q)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("%s", strings.Join(errs, "; "))
	}
	return out, nil
}

func fetchYahooOil(ctx context.Context, symbol, name string) (OilQuote, error) {
	url := fmt.Sprintf(
		"https://query1.finance.yahoo.com/v8/finance/chart/%s?interval=1d&range=1d",
		symbol,
	)
	body, err := httpGet(ctx, url)
	if err != nil {
		return OilQuote{}, fmt.Errorf("%s: %w", symbol, err)
	}

	var resp struct {
		Chart struct {
			Result []struct {
				Meta struct {
					RegularMarketPrice float64 `json:"regularMarketPrice"`
					ChartPreviousClose float64 `json:"chartPreviousClose"`
				} `json:"meta"`
			} `json:"result"`
		} `json:"chart"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return OilQuote{}, fmt.Errorf("%s parse: %w", symbol, err)
	}
	if len(resp.Chart.Result) == 0 {
		return OilQuote{}, fmt.Errorf("%s: empty chart result", symbol)
	}

	meta := resp.Chart.Result[0].Meta
	if meta.RegularMarketPrice == 0 {
		return OilQuote{}, fmt.Errorf("%s: missing market price", symbol)
	}

	change := meta.RegularMarketPrice - meta.ChartPreviousClose
	return OilQuote{
		Name:      name,
		Symbol:    symbol,
		PriceUSD:  meta.RegularMarketPrice,
		ChangeUSD: change,
		Unit:      "USD/thùng",
		Source:    "Yahoo Finance",
	}, nil
}

func parseWebgiaPrice(raw string) int {
	raw = strings.ReplaceAll(strings.TrimSpace(raw), ".", "")
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}
	return n
}

func httpGet(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "golang-news/1.0 (+https://github.com/thanhtoan0306/tony-blogs)")

	client := &http.Client{Timeout: 12 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(io.LimitReader(res.Body, 2<<20))
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("http %d from %s", res.StatusCode, url)
	}
	return body, nil
}
