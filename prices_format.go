package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func formatVND(amount int64) string {
	if amount == 0 {
		return "—"
	}
	s := formatIntDots(amount)
	return s + " ₫"
}

func formatVNDChange(amount int64) string {
	if amount == 0 {
		return "0 ₫"
	}
	prefix := "+"
	if amount < 0 {
		prefix = ""
	}
	return prefix + formatIntDots(amount) + " ₫"
}

func formatVNDPerLiter(amount int) string {
	if amount == 0 {
		return "—"
	}
	return formatIntDots(int64(amount)) + " ₫/lít"
}

func formatUSD(amount float64) string {
	if amount == 0 {
		return "—"
	}
	return fmt.Sprintf("$%.2f", amount)
}

func formatUSDChange(amount float64) string {
	if amount == 0 {
		return "$0.00"
	}
	prefix := "+"
	if amount < 0 {
		prefix = ""
	}
	return prefix + fmt.Sprintf("$%.2f", amount)
}

func formatPriceTime(t time.Time) string {
	if t.IsZero() {
		return "—"
	}
	return t.In(vnLocation()).Format("15:04 · 02/01/2006 MST")
}

func formatNextUpdate(t time.Time) string {
	if t.IsZero() {
		return "—"
	}
	return t.In(vnLocation()).Format("15:04 · 02/01/2006")
}

func vnLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		return time.FixedZone("ICT", 7*3600)
	}
	return loc
}

func formatIntDots(n int64) string {
	neg := n < 0
	if neg {
		n = -n
	}
	s := fmt.Sprintf("%d", n)
	var b strings.Builder
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteRune(r)
	}
	out := b.String()
	if neg {
		out = "-" + out
	}
	return out
}

func changeClass(delta int64) string {
	if delta > 0 {
		return "up"
	}
	if delta < 0 {
		return "down"
	}
	return "flat"
}

func changeClassFloat(delta float64) string {
	if delta > 0 {
		return "up"
	}
	if math.Abs(delta) < 0.005 {
		return "flat"
	}
	if delta < 0 {
		return "down"
	}
	return "flat"
}
