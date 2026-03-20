package main

import (
	"fmt"
	"math/rand"
	"strings"
)

type symbol struct {
	glyph  string
	name   string
	weight int
}

type winTier int

const (
	tierLoss winTier = iota
	tierSmall
	tierMedium
	tierBig
	tierJackpot
)

type winResult struct {
	payout    int
	title     string
	detail    string
	tier      winTier
	highlight [3]bool
}

type spinResult struct {
	stops [3]int
	win   winResult
}

type payRule struct {
	triple int
	pair   int
	title  string
}

var reelSymbols = []symbol{
	{glyph: "🍒", name: "Cherry", weight: 35},
	{glyph: "🔔", name: "Bell", weight: 24},
	{glyph: "⭐", name: "Star", weight: 19},
	{glyph: "💎", name: "Diamond", weight: 14},
	{glyph: "7️⃣", name: "Seven", weight: 8},
}

var payRules = map[string]payRule{
	"Cherry":  {triple: 6, pair: 2, title: "Cherry Bloom"},
	"Bell":    {triple: 10, pair: 3, title: "Bell Chorus"},
	"Star":    {triple: 16, pair: 5, title: "Star Shower"},
	"Diamond": {triple: 28, pair: 8, title: "Diamond Rush"},
	"Seven":   {triple: 60, pair: 12, title: "Jackpot"},
}

func spinReels(rng *rand.Rand, bet int) spinResult {
	var stops [3]int
	for i := range stops {
		stops[i] = weightedStop(rng)
	}

	return spinResult{
		stops: stops,
		win:   evaluateWin(stops, bet),
	}
}

func weightedStop(rng *rand.Rand) int {
	total := 0
	for _, s := range reelSymbols {
		total += s.weight
	}

	roll := rng.Intn(total)
	for i, s := range reelSymbols {
		if roll < s.weight {
			return i
		}
		roll -= s.weight
	}

	return len(reelSymbols) - 1
}

func evaluateWin(stops [3]int, bet int) winResult {
	counts := map[int]int{}
	for _, stop := range stops {
		counts[stop]++
	}

	if stops[0] == stops[1] && stops[1] == stops[2] {
		sym := reelSymbols[stops[0]]
		rule := payRules[sym.name]
		tier := tierBig
		if sym.name == "Seven" {
			tier = tierJackpot
		}

		return winResult{
			payout:    bet * rule.triple,
			title:     rule.title,
			detail:    fmt.Sprintf("%s %s %s pays %dx", sym.glyph, sym.glyph, sym.glyph, rule.triple),
			tier:      tier,
			highlight: [3]bool{true, true, true},
		}
	}

	for idx, count := range counts {
		if count == 2 {
			sym := reelSymbols[idx]
			rule := payRules[sym.name]
			highlight := matchingHighlight(stops, idx)
			tier := tierMedium
			if sym.name == "Seven" || sym.name == "Diamond" {
				tier = tierBig
			}
			return winResult{
				payout:    bet * rule.pair,
				title:     fmt.Sprintf("%s Pair", sym.name),
				detail:    fmt.Sprintf("Two %s symbols pay %dx", strings.ToLower(sym.name), rule.pair),
				tier:      tier,
				highlight: highlight,
			}
		}
	}

	cherryCount := counts[0]
	if cherryCount > 0 {
		multiplier := cherryCount
		title := "Cherry Tap"
		if cherryCount == 2 {
			title = "Cherry Double"
			multiplier = 3
		}
		return winResult{
			payout:    bet * multiplier,
			title:     title,
			detail:    fmt.Sprintf("%d cherry%s on the payline", cherryCount, plural(cherryCount)),
			tier:      tierSmall,
			highlight: matchingHighlight(stops, 0),
		}
	}

	return winResult{
		title:  "Miss",
		detail: "No payout on that spin",
		tier:   tierLoss,
	}
}

func matchingHighlight(stops [3]int, target int) [3]bool {
	var highlight [3]bool
	for i, stop := range stops {
		highlight[i] = stop == target
	}
	return highlight
}

func plural(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}
