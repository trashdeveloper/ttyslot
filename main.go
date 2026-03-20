package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

const frameRate = 75 * time.Millisecond
const (
	initialBalance = 120
	initialBet     = 5
	maxBet         = 25
	historyLimit   = 5
)

type frameMsg time.Time

type reelState struct {
	position        int
	target          int
	stepsRemaining  int
	totalSteps      int
	framesUntilStep int
	settleFrames    int
}

type model struct {
	width        int
	height       int
	rng          *rand.Rand
	themes       []theme
	themeIndex   int
	theme        theme
	styles       styles
	help         help.Model
	keys         keyMap
	reels        [3]reelState
	balance      int
	bet          int
	lastWin      int
	bestWin      int
	spins        int
	totalWagered int
	totalPaid    int
	winStreak    int
	freeSpins    int
	bonusBuys    int
	spinning     bool
	lastResult   spinResult
	status       string
	statusDetail string
	statusTier   winTier
	pulsePhase   int
	history      []spinHistory
}

type spinHistory struct {
	symbols string
	payout  int
}

const (
	freeSpinTriggerOdds = 14
	randomFreeSpinAward = 5
	bonusBuySpins       = 8
	bonusBuyMultiplier  = 18
)

func newModel() model {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	themes := allThemes()
	m := model{
		rng:          rng,
		themes:       themes,
		themeIndex:   0,
		theme:        themes[0],
		help:         help.New(),
		keys:         newKeyMap(),
		balance:      initialBalance,
		bet:          initialBet,
		lastWin:      0,
		status:       "Press Space or Enter to spin the cabinet.",
		statusDetail: "Use the arrow keys to move the bet up or down before you pull.",
		statusTier:   tierSmall,
	}
	m.styles = newStyles(m.theme)

	for i := range m.reels {
		m.reels[i].position = rng.Intn(len(reelSymbols))
		m.reels[i].target = m.reels[i].position
	}

	m.help.ShowAll = false
	return m
}

func tickCmd() tea.Cmd {
	return tea.Tick(frameRate, func(t time.Time) tea.Msg {
		return frameMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case frameMsg:
		m.pulsePhase++
		if m.spinning {
			m.advanceReels()
		}
		return m, tickCmd()
	}

	return m, nil
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.quit):
		return m, tea.Quit
	case key.Matches(msg, m.keys.spin):
		m.startSpin()
	case key.Matches(msg, m.keys.betUp):
		if !m.spinning {
			m.adjustBet(1)
		}
	case key.Matches(msg, m.keys.betDown):
		if !m.spinning {
			m.adjustBet(-1)
		}
	case key.Matches(msg, m.keys.theme):
		if !m.spinning {
			m.cycleTheme()
		}
	case key.Matches(msg, m.keys.reset):
		if !m.spinning {
			m.resetSession()
		}
	case key.Matches(msg, m.keys.buyBonus):
		if !m.spinning {
			m.buyBonus()
		}
	}
	return m, nil
}

func (m *model) startSpin() {
	if m.spinning {
		return
	}
	paidSpin := m.freeSpins == 0
	if paidSpin && m.balance < m.bet {
		m.setStatus("Not enough credits for that bet.", "Drop the wager with the down arrow and try again.", tierLoss)
		return
	}

	if paidSpin {
		m.balance -= m.bet
		m.totalWagered += m.bet
	} else {
		m.freeSpins--
	}
	m.lastWin = 0
	m.lastResult = spinReels(m.rng, m.bet)
	m.spinning = true
	if paidSpin {
		m.setStatus("Spinning...", "The reels are staggering to a stop from left to right.", tierSmall)
	} else {
		m.setStatus("Spinning...", fmt.Sprintf("Free spin in play. %d free spins stay loaded.", m.freeSpins), tierSmall)
	}

	m.prepareReels(m.lastResult.stops)
}

func (m *model) adjustBet(delta int) {
	maxBet := max(1, min(maxBet, max(1, m.balance)))
	m.bet += delta
	if m.bet < 1 {
		m.bet = 1
	}
	if m.bet > maxBet {
		m.bet = maxBet
	}

	m.setStatus(fmt.Sprintf("Bet set to %d credits.", m.bet), "Spin when the cabinet feels right.", tierSmall)
}

func (m *model) advanceReels() {
	allDone := true

	for i := range m.reels {
		reel := &m.reels[i]
		if reel.stepsRemaining > 0 {
			allDone = false
			if reel.framesUntilStep > 0 {
				reel.framesUntilStep--
				continue
			}

			reel.position = (reel.position + 1) % len(reelSymbols)
			reel.stepsRemaining--

			if reel.stepsRemaining == 0 {
				reel.position = reel.target
				reel.settleFrames = 4
			} else {
				reel.framesUntilStep = stepDelay(reel.totalSteps-reel.stepsRemaining, reel.totalSteps)
			}
		}

		if reel.settleFrames > 0 {
			reel.settleFrames--
			allDone = false
		}
	}

	if allDone {
		m.finishSpin()
	}
}

func stepDelay(completed, total int) int {
	if total == 0 {
		return 0
	}

	progress := float64(completed) / float64(total)
	switch {
	case progress < 0.45:
		return 0
	case progress < 0.72:
		return 1
	case progress < 0.88:
		return 2
	default:
		return 3
	}
}

func (m *model) finishSpin() {
	m.spinning = false
	m.spins++
	m.lastWin = m.lastResult.win.payout
	m.balance += m.lastResult.win.payout
	m.totalPaid += m.lastResult.win.payout
	if m.lastWin > m.bestWin {
		m.bestWin = m.lastWin
	}
	if m.lastWin > 0 {
		m.winStreak++
	} else {
		m.winStreak = 0
	}
	m.pushHistory()
	triggered := m.tryAwardFreeSpins()

	if m.lastResult.win.payout == 0 {
		if m.balance == 0 {
			m.setStatus("No payout this round.", "The bankroll is empty. Quit and reopen to start fresh.", tierLoss)
		} else {
			m.setStatus("No payout this round.", "The lights stay warm. Change the bet and spin again.", tierLoss)
		}
		if triggered {
			m.setStatus(
				fmt.Sprintf("Feature trigger! %d free spins", randomFreeSpinAward),
				fmt.Sprintf("Bonus unlocked with %d free spins loaded.", m.freeSpins),
				tierBig,
			)
		}
		return
	}

	m.setStatus(
		fmt.Sprintf("%s! +%d credits", m.lastResult.win.title, m.lastResult.win.payout),
		m.lastResult.win.detail,
		m.lastResult.win.tier,
	)
	if triggered {
		m.setStatus(
			fmt.Sprintf("%s and %d free spins!", m.lastResult.win.title, randomFreeSpinAward),
			fmt.Sprintf("%s Bonus loaded: %d free spins remaining.", m.lastResult.win.detail, m.freeSpins),
			tierBig,
		)
	}
}

func (m *model) cycleTheme() {
	m.themeIndex = (m.themeIndex + 1) % len(m.themes)
	m.theme = m.themes[m.themeIndex]
	m.styles = newStyles(m.theme)
	m.setStatus(fmt.Sprintf("Theme set to %s.", strings.ToLower(m.theme.name)), "Press Space to spin with the new cabinet.", tierSmall)
}

func (m *model) resetSession() {
	m.balance = initialBalance
	m.bet = initialBet
	m.lastWin = 0
	m.bestWin = 0
	m.spins = 0
	m.totalWagered = 0
	m.totalPaid = 0
	m.winStreak = 0
	m.freeSpins = 0
	m.bonusBuys = 0
	m.history = nil
	m.lastResult = spinResult{}
	m.setStatus("Session reset.", "Fresh bankroll loaded. Pull when ready.", tierSmall)
}

func (m *model) pushHistory() {
	entry := spinHistory{
		symbols: m.symbolString(m.lastResult.stops),
		payout:  m.lastResult.win.payout,
	}
	m.history = append([]spinHistory{entry}, m.history...)
	if len(m.history) > historyLimit {
		m.history = m.history[:historyLimit]
	}
}

func (m model) symbolString(stops [3]int) string {
	parts := make([]string, len(stops))
	for i, stop := range stops {
		if stop >= 0 && stop < len(reelSymbols) {
			parts[i] = reelSymbols[stop].glyph
		}
	}
	return strings.Join(parts, " ")
}

func (m model) payoutRate() int {
	if m.totalWagered == 0 {
		return 0
	}
	return int(float64(m.totalPaid) / float64(m.totalWagered) * 100)
}

func (m *model) tryAwardFreeSpins() bool {
	if m.freeSpins > 0 {
		return false
	}
	if m.rng.Intn(freeSpinTriggerOdds) != 0 {
		return false
	}
	m.freeSpins += randomFreeSpinAward
	return true
}

func (m *model) buyBonus() {
	cost := m.bet * bonusBuyMultiplier
	if m.balance < cost {
		m.setStatus("Not enough credits to buy bonus.", fmt.Sprintf("Buy bonus costs %d credits at the current bet.", cost), tierLoss)
		return
	}
	m.balance -= cost
	m.totalWagered += cost
	m.freeSpins += bonusBuySpins
	m.bonusBuys++
	m.setStatus(
		fmt.Sprintf("Bonus bought for %d credits.", cost),
		fmt.Sprintf("%d free spins loaded at %d credits per spin.", bonusBuySpins, m.bet),
		tierBig,
	)
}

func (m *model) setStatus(title, detail string, tier winTier) {
	m.status = title
	m.statusDetail = detail
	m.statusTier = tier
}

func (m *model) prepareReels(targets [3]int) {
	// bleh, the reels still need a little drama.
	for i := range m.reels {
		offset := (targets[i] - m.reels[i].position + len(reelSymbols)) % len(reelSymbols)
		steps := (3+i)*len(reelSymbols) + 8 + offset + m.rng.Intn(4)
		m.reels[i].target = targets[i]
		m.reels[i].stepsRemaining = steps
		m.reels[i].totalSteps = steps
		m.reels[i].framesUntilStep = 0
		m.reels[i].settleFrames = 0
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "ttyslot: %v\n", err)
		os.Exit(1)
	}
}
