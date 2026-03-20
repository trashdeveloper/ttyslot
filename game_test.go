package main

import (
	"math/rand"
	"testing"
)

func TestEvaluateWin(t *testing.T) {
	tests := []struct {
		name   string
		stops  [3]int
		bet    int
		payout int
		tier   winTier
	}{
		{name: "jackpot", stops: [3]int{4, 4, 4}, bet: 5, payout: 300, tier: tierJackpot},
		{name: "triple diamond", stops: [3]int{3, 3, 3}, bet: 2, payout: 56, tier: tierBig},
		{name: "seven pair", stops: [3]int{4, 4, 1}, bet: 3, payout: 36, tier: tierBig},
		{name: "single cherry", stops: [3]int{0, 2, 4}, bet: 4, payout: 4, tier: tierSmall},
		{name: "loss", stops: [3]int{1, 2, 3}, bet: 5, payout: 0, tier: tierLoss},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := evaluateWin(tc.stops, tc.bet)
			if result.payout != tc.payout {
				t.Fatalf("payout = %d, want %d", result.payout, tc.payout)
			}
			if result.tier != tc.tier {
				t.Fatalf("tier = %v, want %v", result.tier, tc.tier)
			}
		})
	}
}

func TestMatchingHighlight(t *testing.T) {
	got := matchingHighlight([3]int{4, 2, 4}, 4)
	want := [3]bool{true, false, true}
	if got != want {
		t.Fatalf("highlight = %v, want %v", got, want)
	}
}

func TestWeightedStopAlwaysInRange(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	for i := 0; i < 200; i++ {
		stop := weightedStop(rng)
		if stop < 0 || stop >= len(reelSymbols) {
			t.Fatalf("stop = %d out of range", stop)
		}
	}
}
