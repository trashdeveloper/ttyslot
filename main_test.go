package main

import (
	"math/rand"
	"strings"
	"testing"
)

func TestPayoutRate(t *testing.T) {
	m := newModel()
	m.totalWagered = 40
	m.totalPaid = 60
	if got := m.payoutRate(); got != 150 {
		t.Fatalf("payoutRate = %d, want 150", got)
	}
}

func TestCycleTheme(t *testing.T) {
	m := newModel()
	start := m.theme.name
	m.cycleTheme()
	if m.theme.name == start {
		t.Fatalf("theme did not change")
	}
}

func TestPushHistoryKeepsNewestFive(t *testing.T) {
	m := newModel()
	for i := 0; i < 7; i++ {
		m.lastResult = spinResult{
			stops: [3]int{0, 1, 2},
			win:   winResult{title: "Test", payout: i},
		}
		m.pushHistory()
	}
	if len(m.history) != 5 {
		t.Fatalf("history length = %d, want 5", len(m.history))
	}
	if m.history[0].payout != 6 {
		t.Fatalf("newest payout = %d, want 6", m.history[0].payout)
	}
}

func TestResetSession(t *testing.T) {
	m := newModel()
	m.balance = 50
	m.bet = 10
	m.bestWin = 100
	m.spins = 9
	m.totalWagered = 40
	m.totalPaid = 30
	m.winStreak = 2
	m.history = []spinHistory{{symbols: "old"}}
	m.resetSession()

	if m.balance != 120 || m.bet != 5 {
		t.Fatalf("reset balance/bet = %d/%d", m.balance, m.bet)
	}
	if m.bestWin != 0 || m.spins != 0 || len(m.history) != 0 {
		t.Fatalf("session stats were not reset")
	}
	if m.freeSpins != 0 || m.bonusBuys != 0 {
		t.Fatalf("bonus state was not reset")
	}
}

func TestBuyBonus(t *testing.T) {
	m := newModel()
	m.balance = 200
	m.bet = 5
	m.buyBonus()
	if m.freeSpins != bonusBuySpins {
		t.Fatalf("freeSpins = %d, want %d", m.freeSpins, bonusBuySpins)
	}
	if m.balance != 200-(5*bonusBuyMultiplier) {
		t.Fatalf("balance = %d after buy", m.balance)
	}
}

func TestTryAwardFreeSpinsOnlyWhenEmpty(t *testing.T) {
	m := newModel()
	m.freeSpins = 2
	if m.tryAwardFreeSpins() {
		t.Fatalf("should not award when spins already loaded")
	}
}

func TestStartSpinUsesFreeSpinWithoutBalance(t *testing.T) {
	m := newModel()
	m.rng = rand.New(rand.NewSource(7))
	m.balance = 0
	m.bet = 5
	m.freeSpins = 2

	m.startSpin()

	if !m.spinning {
		t.Fatalf("expected spin to start with loaded free spins")
	}
	if m.freeSpins != 1 {
		t.Fatalf("freeSpins = %d, want 1", m.freeSpins)
	}
	if m.balance != 0 {
		t.Fatalf("balance = %d, want 0", m.balance)
	}
	if m.totalWagered != 0 {
		t.Fatalf("totalWagered = %d, want 0", m.totalWagered)
	}
}

func TestBuyBonusFailsWithoutCredits(t *testing.T) {
	m := newModel()
	m.balance = 10
	m.bet = 5

	m.buyBonus()

	if m.freeSpins != 0 {
		t.Fatalf("freeSpins = %d, want 0", m.freeSpins)
	}
	if !strings.Contains(strings.ToLower(m.status), "not enough credits") {
		t.Fatalf("status = %q, want insufficient credits message", m.status)
	}
}

func TestViewLargeLayoutContainsCabinetElements(t *testing.T) {
	m := newModel()
	m.width = 120
	m.height = 32

	view := m.View()

	for _, want := range []string{"NEON REELS", "Meters", "Paytable"} {
		if !strings.Contains(view, want) {
			t.Fatalf("view missing %q", want)
		}
	}
	if strings.Contains(view, "Resize to about 64x20") {
		t.Fatalf("large layout rendered compact warning")
	}
}
