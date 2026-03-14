package ui

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"gtree/internal/fs"
)

const minBarWidth = 12

type scanCompleteMsg struct {
	result *fs.ScanResult
	err    error
}

type tickMsg time.Time

type model struct {
	rootPath string
	root     *fs.Node
	focus    *fs.Node
	stack    []*fs.Node
	cursor   int
	width    int
	height   int
	loading  bool
	ticks    int
	err      error
	warnings []string
}

func New(rootPath string) tea.Model {
	return model{
		rootPath: rootPath,
		loading:  true,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(scanTreeCmd(m.rootPath), tickCmd())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tickMsg:
		if m.loading {
			m.ticks++
			return m, tickCmd()
		}
		return m, nil
	case scanCompleteMsg:
		m.loading = false
		m.err = msg.err
		if msg.result != nil {
			m.root = msg.result.Root
			m.focus = msg.result.Root
			m.warnings = msg.result.Warnings
			m.clampCursor()
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		if m.loading || m.err != nil || m.focus == nil {
			return m, nil
		}

		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.focus.Children)-1 {
				m.cursor++
			}
		case "home", "g":
			m.cursor = 0
		case "end", "G":
			if len(m.focus.Children) > 0 {
				m.cursor = len(m.focus.Children) - 1
			}
		case "enter", "right", "l":
			child := m.selected()
			if child != nil && child.IsDir {
				m.stack = append(m.stack, m.focus)
				m.focus = child
				m.cursor = 0
			}
		case "backspace", "left", "h":
			if len(m.stack) > 0 {
				m.focus = m.stack[len(m.stack)-1]
				m.stack = m.stack[:len(m.stack)-1]
				m.cursor = 0
			}
		}

		m.clampCursor()
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	if m.loading {
		return m.renderLoading()
	}

	if m.err != nil {
		return m.renderError()
	}

	if m.focus == nil {
		return "No data available."
	}

	return m.renderTree()
}

func scanTreeCmd(root string) tea.Cmd {
	return func() tea.Msg {
		result, err := fs.Scan(root)
		return scanCompleteMsg{result: result, err: err}
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(140*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) selected() *fs.Node {
	if m.focus == nil || len(m.focus.Children) == 0 {
		return nil
	}

	if m.cursor < 0 || m.cursor >= len(m.focus.Children) {
		return nil
	}

	return m.focus.Children[m.cursor]
}

func (m *model) clampCursor() {
	if m.focus == nil || len(m.focus.Children) == 0 {
		m.cursor = 0
		return
	}

	if m.cursor < 0 {
		m.cursor = 0
	}

	if m.cursor >= len(m.focus.Children) {
		m.cursor = len(m.focus.Children) - 1
	}
}

func (m model) renderLoading() string {
	dots := strings.Repeat(".", m.ticks%4)
	title := titleStyle.Render("gtree")
	message := bodyStyle.Render(fmt.Sprintf("Scanning %s%s", m.rootPath, dots))
	help := mutedStyle.Render("Recursive size crawl in progress")
	return lipgloss.JoinVertical(lipgloss.Left, title, "", message, help)
}

func (m model) renderError() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("gtree"),
		"",
		errorStyle.Render(fmt.Sprintf("Scan failed: %v", m.err)),
		mutedStyle.Render("Press q to quit."),
	)
}

func (m model) renderTree() string {
	header := m.renderHeader()
	rows := m.renderRows()
	detail := m.renderDetail()
	help := m.renderHelp()

	return lipgloss.JoinVertical(lipgloss.Left, header, "", rows, "", detail, "", help)
}

func (m model) renderHeader() string {
	breadcrumbs := []string{filepath.Base(m.root.Path)}
	for _, node := range m.stack {
		if node.Path == m.root.Path {
			continue
		}
		breadcrumbs = append(breadcrumbs, node.Name)
	}
	if m.focus.Path != m.root.Path {
		breadcrumbs = append(breadcrumbs, m.focus.Name)
	}

	pathLine := mutedStyle.Render(strings.Join(breadcrumbs, " / "))
	summary := fmt.Sprintf(
		"%s total  |  %d entries  |  %d warnings",
		humanizeBytes(m.focus.Size),
		len(m.focus.Children),
		len(m.warnings),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("gtree"),
		pathLine,
		bodyStyle.Render(summary),
	)
}

func (m model) renderRows() string {
	if len(m.focus.Children) == 0 {
		return emptyStyle.Render("This folder is empty. Use backspace/h/left to return.")
	}

	available := m.height - 10
	if available < 5 {
		available = 5
	}

	start := 0
	if m.cursor >= available {
		start = m.cursor - available + 1
	}
	end := start + available
	if end > len(m.focus.Children) {
		end = len(m.focus.Children)
	}

	rows := make([]string, 0, end-start)
	for index := start; index < end; index++ {
		rows = append(rows, m.renderRow(index, m.focus.Children[index]))
	}

	return strings.Join(rows, "\n")
}

func (m model) renderRow(index int, node *fs.Node) string {
	selected := index == m.cursor
	marker := " "
	if selected {
		marker = ">"
	}

	name := node.Name
	if node.IsDir {
		name += "/"
	}

	share := 0.0
	if m.focus.Size > 0 {
		share = float64(node.Size) / float64(m.focus.Size)
	}

	barWidth := m.width - 40
	if barWidth < minBarWidth {
		barWidth = minBarWidth
	}

	bar := makeBar(barWidth, share, selected)
	row := fmt.Sprintf(
		"%s %-28s %9s %6.1f%% %s",
		marker,
		truncate(name, 28),
		humanizeBytes(node.Size),
		share*100,
		bar,
	)

	if selected {
		return selectedRowStyle.Render(row)
	}

	return rowStyle.Render(row)
}

func (m model) renderDetail() string {
	selected := m.selected()
	if selected == nil {
		return mutedStyle.Render("No item selected.")
	}

	kind := "file"
	action := "Press enter on a folder to drill in."
	if selected.IsDir {
		kind = "folder"
		action = "Press enter to focus this folder."
	}

	share := 0.0
	if m.focus.Size > 0 {
		share = float64(selected.Size) / float64(m.focus.Size)
	}

	lines := []string{
		bodyStyle.Render(fmt.Sprintf("Selected: %s", selected.Path)),
		mutedStyle.Render(fmt.Sprintf("%s | %s | %.1f%% of current view", kind, humanizeBytes(selected.Size), share*100)),
		mutedStyle.Render(action),
	}

	if len(m.stack) > 0 {
		lines = append(lines, mutedStyle.Render("Backspace, h, or left arrow returns to the parent view."))
	}

	return strings.Join(lines, "\n")
}

func (m model) renderHelp() string {
	controls := "j/k or arrows: move   enter: open folder   h/backspace/left: back   g/G: top/bottom   q: quit"
	return mutedStyle.Render(controls)
}

func makeBar(width int, share float64, selected bool) string {
	filled := int(share * float64(width))
	if share > 0 && filled == 0 {
		filled = 1
	}
	if filled > width {
		filled = width
	}

	fillChar := "█"
	emptyChar := "░"
	if selected {
		fillChar = "▓"
	}

	return strings.Repeat(fillChar, filled) + strings.Repeat(emptyChar, width-filled)
}

func humanizeBytes(size int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	value := float64(size)
	unit := 0

	for value >= 1024 && unit < len(units)-1 {
		value /= 1024
		unit++
	}

	if unit == 0 {
		return fmt.Sprintf("%d %s", size, units[unit])
	}

	return fmt.Sprintf("%.1f %s", value, units[unit])
}

func truncate(value string, width int) string {
	runes := []rune(value)
	if len(runes) <= width {
		return value
	}

	if width <= 1 {
		return string(runes[:width])
	}

	if width <= 3 {
		return string(runes[:width])
	}

	return string(runes[:width-3]) + "..."
}
