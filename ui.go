package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

type theme struct {
	name       string
	text       string
	muted      string
	border     string
	accent     string
	accentSoft string
	success    string
	danger     string
	cherry     string
	bell       string
	star       string
	diamond    string
	seven      string
}

type styles struct {
	app           lipgloss.Style
	compact       lipgloss.Style
	footer        lipgloss.Style
	marquee       lipgloss.Style
	marqueeTitle  lipgloss.Style
	lampOn        lipgloss.Style
	lampOff       lipgloss.Style
	cabinet       lipgloss.Style
	rail          lipgloss.Style
	railTitle     lipgloss.Style
	railText      lipgloss.Style
	railMuted     lipgloss.Style
	glass         lipgloss.Style
	stage         lipgloss.Style
	badge         lipgloss.Style
	ticker        lipgloss.Style
	reel          lipgloss.Style
	reelBlur      lipgloss.Style
	reelCenter    lipgloss.Style
	reelSpin      lipgloss.Style
	reelWin       lipgloss.Style
	reelJackpot   lipgloss.Style
	arrow         lipgloss.Style
	statusText    lipgloss.Style
	statusMuted   lipgloss.Style
	statusLoss    lipgloss.Style
	statusWin     lipgloss.Style
	statusJackpot lipgloss.Style
	meterLabel    lipgloss.Style
}

type layoutMode int

const (
	layoutCompact layoutMode = iota
	layoutSmall
	layoutMedium
	layoutLarge
)

type layoutSpec struct {
	mode        layoutMode
	centerWidth int
	reelWidth   int
	lampCount   int
	showHandle  bool
	showFooter  bool
}

type keyMap struct {
	spin     key.Binding
	betUp    key.Binding
	betDown  key.Binding
	theme    key.Binding
	buyBonus key.Binding
	reset    key.Binding
	quit     key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		spin: key.NewBinding(
			key.WithKeys(" ", "enter"),
			key.WithHelp("space/enter", "spin"),
		),
		betUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("up", "bet +"),
		),
		betDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("down", "bet -"),
		),
		theme: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "theme"),
		),
		buyBonus: key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "buy bonus"),
		),
		reset: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reset"),
		),
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

func defaultTheme() theme {
	return theme{
		name:       "Nord",
		text:       "#ECEFF4",
		muted:      "#A3BE8C",
		border:     "#4C566A",
		accent:     "#88C0D0",
		accentSoft: "#BF616A",
		success:    "#A3BE8C",
		danger:     "#BF616A",
		cherry:     "#BF616A",
		bell:       "#EBCB8B",
		star:       "#EBCB8B",
		diamond:    "#81A1C1",
		seven:      "#D08770",
	}
}

func allThemes() []theme {
	return []theme{
		defaultTheme(),
		{
			name:       "Gruvbox",
			text:       "#EBDBB2",
			muted:      "#BDAE93",
			border:     "#665C54",
			accent:     "#FABD2F",
			accentSoft: "#FB4934",
			success:    "#B8BB26",
			danger:     "#FB4934",
			cherry:     "#FB4934",
			bell:       "#FABD2F",
			star:       "#FE8019",
			diamond:    "#83A598",
			seven:      "#D3869B",
		},
		{
			name:       "Catppuccin",
			text:       "#CDD6F4",
			muted:      "#A6ADC8",
			border:     "#585B70",
			accent:     "#F9E2AF",
			accentSoft: "#F38BA8",
			success:    "#A6E3A1",
			danger:     "#F38BA8",
			cherry:     "#F38BA8",
			bell:       "#F9E2AF",
			star:       "#FAB387",
			diamond:    "#89DCEB",
			seven:      "#CBA6F7",
		},
		{
			name:       "Tokyo Night",
			text:       "#C0CAF5",
			muted:      "#9AA5CE",
			border:     "#414868",
			accent:     "#E0AF68",
			accentSoft: "#F7768E",
			success:    "#9ECE6A",
			danger:     "#F7768E",
			cherry:     "#F7768E",
			bell:       "#E0AF68",
			star:       "#FF9E64",
			diamond:    "#7DCFFF",
			seven:      "#BB9AF7",
		},
	}
}

func newStyles(p theme) styles {
	return styles{
		app: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.text)),
		compact: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.accent)).
			Bold(true),
		footer: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.muted)).
			PaddingTop(0),
		marquee: lipgloss.NewStyle().
			Padding(0, 2).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(p.border)),
		marqueeTitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.accentSoft)).
			Bold(true),
		lampOn: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.accent)).
			Bold(true),
		lampOff: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A442E")),
		cabinet: lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color(p.border)),
		rail: lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(p.border)),
		railTitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.accent)).
			Bold(true),
		railText: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.text)),
		railMuted: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.muted)),
		glass: lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color(p.border)),
		stage: lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(p.border)),
		badge: lipgloss.NewStyle().
			Padding(0, 2).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(p.border)).
			Foreground(lipgloss.Color(p.accent)).
			Bold(true),
		ticker: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.muted)),
		reel: lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(p.border)),
		reelBlur: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.muted)).
			Align(lipgloss.Center),
		reelCenter: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.text)).
			Align(lipgloss.Center).
			Bold(true),
		reelSpin: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.accent)).
			Align(lipgloss.Center).
			Bold(true),
		reelWin: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.success)).
			Align(lipgloss.Center).
			Bold(true).
			Underline(true),
		reelJackpot: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.danger)).
			Align(lipgloss.Center).
			Bold(true).
			Underline(true),
		arrow: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.accent)).
			Bold(true),
		statusText: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.text)).
			Bold(true),
		statusMuted: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.muted)),
		statusLoss: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.accentSoft)).
			Bold(true),
		statusWin: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.success)).
			Bold(true),
		statusJackpot: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.danger)).
			Bold(true),
		meterLabel: lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.accent)),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.spin, k.betUp, k.betDown, k.buyBonus, k.theme, k.quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.spin, k.betUp, k.betDown},
		{k.theme, k.buyBonus, k.reset},
		{k.quit},
	}
}

func (m model) View() string {
	layout := m.layoutSpec()
	if layout.mode == layoutCompact {
		msg := m.styles.compact.Render("Resize to about 64x20 for the slot cabinet.")
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, msg)
	}

	machine := m.renderMachine(layout)
	footer := ""
	if layout.showFooter {
		footer = m.styles.footer.Render(m.help.View(m.keys))
	}

	content := lipgloss.JoinVertical(lipgloss.Center, machine, footer)
	vAlign := lipgloss.Center
	if layout.mode == layoutSmall || lipgloss.Height(content) >= m.height-1 {
		vAlign = lipgloss.Top
	}

	return m.styles.app.Render(lipgloss.Place(m.width, m.height, lipgloss.Center, vAlign, content))
}

func (m model) layoutSpec() layoutSpec {
	switch {
	case m.width >= 120 && m.height >= 30:
		return layoutSpec{
			mode:        layoutLarge,
			centerWidth: 60,
			reelWidth:   10,
			lampCount:   8,
			showHandle:  true,
			showFooter:  true,
		}
	case m.width >= 96 && m.height >= 24:
		return layoutSpec{
			mode:        layoutMedium,
			centerWidth: 48,
			reelWidth:   8,
			lampCount:   6,
			showHandle:  false,
			showFooter:  true,
		}
	case m.width >= 64 && m.height >= 20:
		return layoutSpec{
			mode:        layoutSmall,
			centerWidth: min(50, m.width-8),
			reelWidth:   7,
			lampCount:   4,
			showFooter:  false,
		}
	default:
		return layoutSpec{mode: layoutCompact}
	}
}

func (m model) renderMachine(layout layoutSpec) string {
	var cabinet string
	if layout.mode == layoutSmall {
		cabinet = m.renderCompactCabinet(layout)
		marquee := m.renderMarquee(layout, lipgloss.Width(cabinet))
		return lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.NewStyle().Width(lipgloss.Width(cabinet)).Align(lipgloss.Center).Render(marquee),
			cabinet,
		)
	} else {
		cabinet = m.renderCabinet(layout)
	}

	marquee := m.renderMarquee(layout, lipgloss.Width(cabinet))
	stack := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().Width(lipgloss.Width(cabinet)).Align(lipgloss.Center).Render(marquee),
		cabinet,
	)
	if !layout.showHandle {
		return stack
	}

	handle := lipgloss.NewStyle().
		MarginTop(lipgloss.Height(marquee) + 2).
		Render(m.renderHandle())
	return lipgloss.JoinHorizontal(lipgloss.Top, stack, " ", handle)
}

func (m model) renderMarquee(layout layoutSpec, targetWidth int) string {
	innerWidth := max(24, targetWidth-m.styles.marquee.GetHorizontalFrameSize())
	sideLampCount := max(2, min(layout.lampCount, (innerWidth-18)/4))
	lamps := m.renderLampRow(sideLampCount)

	title := lipgloss.NewStyle().
		Width(innerWidth).
		Align(lipgloss.Center).
		Render(
			lamps +
				"  " +
				m.styles.marqueeTitle.Render("NEON REELS") +
				"  " +
				lamps,
		)

	return m.styles.marquee.
		Width(innerWidth).
		Render(title)
}

func (m model) renderCabinet(layout layoutSpec) string {
	slot := m.renderSlotWindow(layout)
	slotWidth := lipgloss.Width(slot)
	stack := lipgloss.JoinVertical(
		lipgloss.Center,
		slot,
		m.renderLowerDeck(slotWidth),
	)
	return m.styles.cabinet.Render(stack)
}

func (m model) renderCompactCabinet(layout layoutSpec) string {
	width := layout.centerWidth + 4
	stack := lipgloss.JoinVertical(
		lipgloss.Center,
		m.renderSlotWindow(layout),
		m.renderCompactStatus(width),
		m.renderCompactPaytable(width),
	)
	return m.styles.cabinet.Render(stack)
}

func (m model) renderLowerDeck(totalWidth int) string {
	panelWidth := max(28, (totalWidth-6)/2)
	panelHeight := 8
	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.renderMetersRail(panelWidth, panelHeight),
		"  ",
		m.renderStatusPayDeck(panelWidth, panelHeight),
	)
	return lipgloss.NewStyle().
		Width(totalWidth).
		Align(lipgloss.Center).
		Render(row)
}

func (m model) renderMetersRail(width, height int) string {
	bodyWidth := panelBodyWidth(width)
	bankrollBar := m.renderProgressBar(bodyWidth, "bank", float64(min(m.balance, 200))/200, m.theme.success)
	betBar := m.renderProgressBar(bodyWidth, "bet", float64(m.bet)/25, m.theme.accent)
	heatValue := 0.12
	if m.spinning {
		heatValue = 0.72
	} else if m.lastWin > 0 {
		heatValue = float64(min(100, m.lastWin*2)) / 100
	}
	heatBar := m.renderProgressBar(bodyWidth, "heat", heatValue, m.theme.danger)

	lines := []string{
		m.styles.railTitle.Render("Meters"),
		m.renderStatLine("balance", fmt.Sprintf("%d", m.balance), width),
		m.renderStatLine("bet", fmt.Sprintf("%d", m.bet), width),
		m.renderStatLine("last win", fmt.Sprintf("%d", m.lastWin), width),
		m.renderStatLine("spins", fmt.Sprintf("%d", m.spins), width),
		m.renderStatLine("best win", fmt.Sprintf("%d", m.bestWin), width),
		m.renderStatLine("rtp", fmt.Sprintf("%d%%", m.payoutRate()), width),
		bankrollBar,
		betBar,
		heatBar,
	}
	if m.freeSpins > 0 {
		lines = append(lines, m.styles.statusWin.Render(m.fitText(fmt.Sprintf("free spins loaded: %d", m.freeSpins), bodyWidth)))
	}

	return m.renderRail(width, height, lines...)
}

func (m model) renderStatusPayDeck(width, height int) string {
	lines := []string{
		m.styles.railTitle.Render("Paytable"),
		m.renderPayLine("7️⃣ 7️⃣ 7️⃣", "60x", width, m.theme.seven),
		m.renderPayLine("💎 💎 💎", "28x", width, m.theme.diamond),
		m.renderPayLine("⭐ ⭐ ⭐", "16x", width, m.theme.star),
		m.renderPayLine("🔔 🔔 🔔", "10x", width, m.theme.bell),
		m.renderPayLine("🍒 anywhere", "1x+", width, m.theme.cherry),
		m.currentStatusStyle().Render(m.fitText(strings.ToLower(m.status), panelBodyWidth(width))),
		m.styles.statusMuted.Render(m.fitText(strings.ToLower(m.statusDetail), panelBodyWidth(width))),
		m.styles.railMuted.Render(m.fitText("recent: "+m.renderRecentHistory(width-10), panelBodyWidth(width))),
	}
	if m.freeSpins > 0 {
		lines = append(lines, m.styles.statusWin.Render(m.fitText(fmt.Sprintf("%d free spins ready", m.freeSpins), panelBodyWidth(width))))
	} else {
		lines = append(lines, m.styles.railMuted.Render(m.fitText(fmt.Sprintf("buy bonus: %d credits", m.bet*bonusBuyMultiplier), panelBodyWidth(width))))
	}
	return m.renderRail(width, height, lines...)
}

func (m model) renderRail(width, height int, lines ...string) string {
	body := lipgloss.JoinVertical(lipgloss.Center, m.centerLines(width, lines...)...)
	return m.styles.rail.
		Width(width).
		Height(height).
		Render(body)
}

func (m model) renderSlotWindow(layout layoutSpec) string {
	lampRow := m.renderLampRow(max(4, layout.lampCount))

	titleLine := lipgloss.NewStyle().
		Width(layout.centerWidth).
		Align(lipgloss.Center).
		Render(lampRow + "  " + m.styles.badge.Render("neon reels") + "  " + lampRow)

	reels := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.renderReelColumn(0, layout),
		"  ",
		m.renderReelColumn(1, layout),
		"  ",
		m.renderReelColumn(2, layout),
	)
	reelBlockWidth := lipgloss.Width(reels)

	payline := lipgloss.JoinHorizontal(
		lipgloss.Center,
		m.styles.arrow.Render(m.animatedArrow(true)),
		"  ",
		lipgloss.NewStyle().Width(reelBlockWidth).Align(lipgloss.Center).Render(reels),
		"  ",
		m.styles.arrow.Render(m.animatedArrow(false)),
	)

	payline = lipgloss.NewStyle().
		Width(layout.centerWidth).
		Align(lipgloss.Center).
		Render(payline)

	ticker := lipgloss.NewStyle().
		Width(layout.centerWidth).
		Align(lipgloss.Center).
		Render(m.renderStageTicker(layout.centerWidth))

	stage := m.styles.stage.
		Width(layout.centerWidth).
		Render(lipgloss.JoinVertical(lipgloss.Center, titleLine, payline, ticker))

	return m.styles.glass.Render(stage)
}

func (m model) renderReelColumn(idx int, layout layoutSpec) string {
	reel := m.reels[idx]
	count := len(reelSymbols)
	window := [3]symbol{
		reelSymbols[(reel.position-1+count)%count],
		reelSymbols[reel.position],
		reelSymbols[(reel.position+1)%count],
	}

	rows := make([]string, 0, 3)
	for rowIdx, sym := range window {
		style := m.styles.reelBlur.Width(layout.reelWidth)
		if rowIdx == 1 {
			style = m.styles.reelCenter.Width(layout.reelWidth)
		}
		if m.spinning && reel.stepsRemaining > 0 {
			style = m.styles.reelSpin.Width(layout.reelWidth)
		}
		if rowIdx == 1 && m.lastResult.win.payout > 0 && m.lastResult.win.highlight[idx] && !m.spinning {
			if m.lastResult.win.tier == tierJackpot {
				style = m.styles.reelJackpot.Width(layout.reelWidth)
			} else {
				style = m.styles.reelWin.Width(layout.reelWidth)
			}
		}
		rows = append(rows, style.Foreground(lipgloss.Color(m.symbolColor(sym.name))).Render(sym.glyph))
	}

	return m.styles.reel.Render(lipgloss.JoinVertical(lipgloss.Center, rows...))
}

func (m model) renderCompactStatus(width int) string {
	lines := []string{
		m.styles.railTitle.Render("Machine"),
		m.currentStatusStyle().Render(m.fitText(strings.ToLower(m.status), width)),
		m.styles.statusMuted.Render(m.fitText(strings.ToLower(m.statusDetail), width)),
		m.styles.railMuted.Render(m.fitText(fmt.Sprintf("balance %d  bet %d  win %d", m.balance, m.bet, m.lastWin), width)),
	}
	return m.styles.rail.Width(width).Render(lipgloss.JoinVertical(lipgloss.Center, m.centerLines(width, lines...)...))
}

func (m model) renderCompactPaytable(width int) string {
	lines := []string{
		m.styles.railTitle.Render("Paytable"),
		m.styles.railText.Render(m.fitText("7️⃣ 7️⃣ 7️⃣ 60x  ·  💎 💎 💎 28x", width)),
		m.styles.railText.Render(m.fitText("⭐ ⭐ ⭐ 16x  ·  🔔 🔔 🔔 10x", width)),
		m.styles.railText.Render(m.fitText("🍒 anywhere 1x+", width)),
	}
	return m.styles.rail.Width(width).Render(lipgloss.JoinVertical(lipgloss.Center, m.centerLines(width, lines...)...))
}

func (m model) renderRecentHistory(width int) string {
	if len(m.history) == 0 {
		return m.fitText("spin to build history", width)
	}
	entry := m.history[0]
	line := entry.symbols
	if entry.payout > 0 {
		line += fmt.Sprintf("  +%d", entry.payout)
	} else {
		line += "  miss"
	}
	return m.fitText(strings.ToLower(line), width)
}

func (m model) renderStageTicker(width int) string {
	line := m.styles.ticker.Render("free spins can trigger randomly  ·  buy bonus scales with bet")
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(m.fitText(line, width))
}

func (m model) centerLines(width int, lines ...string) []string {
	// yep, this one just wants to sit in the middle.
	centered := make([]string, 0, len(lines))
	bodyWidth := panelBodyWidth(width)
	for _, line := range lines {
		centered = append(centered, lipgloss.NewStyle().Width(bodyWidth).Align(lipgloss.Center).Render(line))
	}
	return centered
}

func (m model) renderStatLine(label, value string, width int) string {
	line := m.styles.railMuted.Render(label) + "  " + m.styles.railText.Render(value)
	bodyWidth := panelBodyWidth(width)
	return lipgloss.NewStyle().Width(bodyWidth).Align(lipgloss.Center).Render(m.fitText(line, bodyWidth))
}

func (m model) renderPayLine(label, payout string, width int, color string) string {
	line := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Render(label) + "  " + m.styles.railText.Render(payout)
	bodyWidth := panelBodyWidth(width)
	return lipgloss.NewStyle().Width(bodyWidth).Align(lipgloss.Center).Render(m.fitText(line, bodyWidth))
}

func (m model) renderProgressBar(width int, label string, value float64, color string) string {
	if width < 14 {
		width = 14
	}
	labelWidth := max(4, min(6, width/3))
	barWidth := max(8, width-labelWidth-1)
	bar := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(barWidth),
		progress.WithoutPercentage(),
		progress.WithScaledGradient(color, color),
	)

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		m.styles.meterLabel.Width(labelWidth).Render(strings.ToLower(label)),
		bar.ViewAs(value),
	)
}

func (m model) renderLampRow(count int) string {
	bulbs := make([]string, count)
	phase := (m.pulsePhase / 2) % max(1, count)
	for i := range bulbs {
		active := (i+phase)%2 == 0
		if m.spinning {
			active = (i+phase)%3 != 0
		}
		if active {
			bulbs[i] = m.styles.lampOn.Render("●")
		} else {
			bulbs[i] = m.styles.lampOff.Render("●")
		}
	}
	return strings.Join(bulbs, " ")
}

func (m model) animatedArrow(left bool) string {
	if left {
		return "◀"
	}
	return "▶"
}

func (m model) symbolColor(name string) string {
	switch name {
	case "Cherry":
		return m.theme.cherry
	case "Bell":
		return m.theme.bell
	case "Star":
		return m.theme.star
	case "Diamond":
		return m.theme.diamond
	case "Seven":
		return m.theme.seven
	default:
		return m.theme.text
	}
}

func (m model) currentStatusStyle() lipgloss.Style {
	switch m.statusTier {
	case tierLoss:
		return m.styles.statusLoss
	case tierSmall, tierMedium, tierBig:
		return m.styles.statusWin
	case tierJackpot:
		return m.styles.statusJackpot
	default:
		return m.styles.statusText
	}
}

func (m model) fitText(s string, width int) string {
	if width <= 3 {
		return s
	}
	if lipgloss.Width(s) <= width {
		return s
	}

	var b strings.Builder
	current := 0
	for _, r := range s {
		rw := runeWidth(r)
		if current+rw > width-1 {
			break
		}
		b.WriteRune(r)
		current += rw
	}
	return strings.TrimSpace(b.String()) + "…"
}

func runeWidth(r rune) int {
	if r == 0 {
		return 0
	}
	if r < utf8.RuneSelf {
		return 1
	}
	return lipgloss.Width(string(r))
}

func panelBodyWidth(width int) int {
	return max(1, width-2)
}

func (m model) renderHandle() string {
	shaft := lipgloss.NewStyle().Foreground(lipgloss.Color(m.theme.border)).Bold(true)
	ballColor := m.theme.accent
	if m.spinning {
		ballColor = m.theme.danger
	}
	ball := lipgloss.NewStyle().Foreground(lipgloss.Color(ballColor)).Bold(true).Render("◉")
	pivot := lipgloss.NewStyle().Foreground(lipgloss.Color(m.theme.border)).Bold(true).Render("◎")

	if m.spinning {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			"      ",
			"      ",
			pivot+shaft.Render("━━╮"),
			"   "+shaft.Render("╲"),
			"    "+shaft.Render("╲"),
			"     "+ball,
		)
	}

	if m.anyReelSettling() {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			"      ",
			"    "+ball,
			"   "+shaft.Render("╱"),
			pivot+shaft.Render("─╯"),
			"  "+shaft.Render("╲"),
			"   "+shaft.Render("╲"),
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		"    "+ball,
		"   "+shaft.Render("╱"),
		"  "+shaft.Render("╱"),
		pivot+shaft.Render("╯"),
		" "+shaft.Render("│"),
		" "+shaft.Render("│"),
	)
}

func (m model) anyReelSettling() bool {
	for _, reel := range m.reels {
		if reel.settleFrames > 0 {
			return true
		}
	}
	return false
}
